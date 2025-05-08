package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/plaenkler/ddns-updater/pkg/config"
)

const hetznerBaseURL = "https://dns.hetzner.com/api/v1"

type client struct {
	APIToken string
}

type UpdateHetznerRequest struct {
	APIToken     string
	RecordZoneID string
	RecordType   string
	RecordName   string
	RecordTTL    uint32 `json:",string"`
}

type Record struct {
	ID     string `json:"id,omitempty"`
	ZoneID string `json:"zone_id,omitempty"`
	Type   string `json:"type,omitempty"`
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	TTL    uint32 `json:"ttl,omitempty"`
	Error  string `json:"error,omitempty"`
}

func UpdateHetzner(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateHetznerRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	c := &client{
		APIToken: r.APIToken,
	}

	rec, found, err := c.findRecord(r.RecordZoneID, r.RecordName)
	if err != nil {
		return fmt.Errorf("failed to find record: %w", err)
	}

	if !found {
		rec := &Record{
			ZoneID: r.RecordZoneID,
			Type:   r.RecordType,
			Name:   r.RecordName,
			TTL:    r.RecordTTL,
			Value:  ipAddr,
		}

		err = c.createRecord(rec)
		if err != nil {
			return fmt.Errorf("failed to create record: %w", err)
		}
		return nil
	}

	rec.Value = ipAddr
	rec.TTL = r.RecordTTL
	err = c.updateRecord(rec)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}
	return nil
}

func (c *client) updateRecord(record *Record) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(hetznerBaseURL+"/records/%s", record.ID)
	body, err := c.fetch(http.MethodPut, url, data)
	if err != nil {
		return err
	}

	var rec Record
	if err := json.Unmarshal(body, &rec); err != nil {
		return err
	}

	if rec.Error != "" {
		return fmt.Errorf("failed to update record: %s", rec.Error)
	}

	return nil
}

func (c *client) createRecord(record *Record) error {
	if record.ZoneID == "" {
		return fmt.Errorf("zone ID is required")
	}

	if record.Type == "" {
		return fmt.Errorf("record type is required")
	}

	if record.Name == "" {
		return fmt.Errorf("record name is required")
	}

	if record.Value == "" {
		return fmt.Errorf("record value is required")
	}

	if record.TTL == 0 {
		record.TTL = config.Get().Interval
	}

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(hetznerBaseURL + "/records")
	body, err := c.fetch(http.MethodPost, url, data)
	if err != nil {
		return err
	}
	var rec Record
	if err := json.Unmarshal(body, &rec); err != nil {
		return err
	}

	if rec.Error != "" {
		return fmt.Errorf("failed to create record: %s", rec.Error)
	}

	return nil
}

func (c *client) findRecord(zoneID, recordName string) (*Record, bool, error) {
	url := fmt.Sprintf(hetznerBaseURL+"/records?zone_id=%s&name=%s", zoneID, recordName)
	body, err := c.fetch(http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err
	}
	var records struct {
		Records []Record `json:"records"`
		Error   string   `json:"error"`
	}
	if err := json.Unmarshal(body, &records); err != nil {
		return nil, false, err
	}
	if records.Error != "" {
		return nil, false, fmt.Errorf("failed to find record: %s", records.Error)
	}
	if len(records.Records) == 0 {
		return nil, false, nil
	}
	return &records.Records[0], true, nil
}

func (c *client) fetch(method string, url string, body []byte) ([]byte, error) {
	req, _ := http.NewRequest(method, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Auth-API-Token", c.APIToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch: %s", resp.Status)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
