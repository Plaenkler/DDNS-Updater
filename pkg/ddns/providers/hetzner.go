package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/plaenkler/ddns-updater/pkg/config"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
)

const (
	hetznerBaseURL        = "https://dns.hetzner.com/api/v1/records"
	contentType           = "Content-Type"
	contentTypeJSON       = "application/json"
	headerAuthToken       = "Auth-API-Token"
	errInvalidRequestType = "invalid request type: %T"
	errAPICreate          = "failed to create record: %s"
	errAPIUpdate          = "failed to update record: %s"
	errAPIFind            = "failed to find record: %s"
	errAPIStatus          = "API returned status: %s"
)

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
	ID     string `json:"id"`
	ZoneID string `json:"zone_id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	TTL    uint32 `json:"ttl"`
	Error  string `json:"error"`
}

func UpdateHetzner(request any, ipAddr string) error {
	r, ok := request.(*UpdateHetznerRequest)
	if !ok {
		return fmt.Errorf(errInvalidRequestType, request)
	}
	c := &client{APIToken: r.APIToken}
	rec, found, err := c.findRecord(r.RecordZoneID, r.RecordName)
	if err != nil {
		return fmt.Errorf("find record: %w", err)
	}
	if !found {
		newRecord := &Record{
			ZoneID: r.RecordZoneID,
			Type:   r.RecordType,
			Name:   r.RecordName,
			TTL:    r.RecordTTL,
			Value:  ipAddr,
		}
		return c.createRecord(newRecord)
	}
	rec.Value = ipAddr
	rec.TTL = r.RecordTTL
	return c.updateRecord(rec)
}

func (c *client) updateRecord(record *Record) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal update: %w", err)
	}
	url := hetznerBaseURL + "/" + record.ID
	body, err := c.fetch(http.MethodPut, url, data)
	if err != nil {
		return err
	}
	return c.handleAPIResponse(body, "update")
}

func (c *client) createRecord(record *Record) error {
	if err := c.validateRecord(record); err != nil {
		return err
	}
	if record.TTL == 0 {
		record.TTL = config.Get().Interval
	}
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal create: %w", err)
	}
	body, err := c.fetch(http.MethodPost, hetznerBaseURL, data)
	if err != nil {
		return err
	}
	return c.handleAPIResponse(body, "create")
}

func (c *client) findRecord(zoneID, recordName string) (*Record, bool, error) {
	url := fmt.Sprintf("%s?zone_id=%s&name=%s", hetznerBaseURL, zoneID, recordName)
	body, err := c.fetch(http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err
	}
	var res struct {
		Records []Record `json:"records"`
		Error   string   `json:"error"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, false, fmt.Errorf("unmarshal find: %w", err)
	}
	if res.Error != "" {
		return nil, false, fmt.Errorf(errAPIFind, res.Error)
	}
	if len(res.Records) == 0 {
		return nil, false, nil
	}
	return &res.Records[0], true, nil
}

func (c *client) fetch(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set(contentType, contentTypeJSON)
	req.Header.Set(headerAuthToken, c.APIToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer log.ErrorClose(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errAPIStatus, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func (c *client) validateRecord(r *Record) error {
	switch {
	case r.ZoneID == "":
		return fmt.Errorf("zone ID is required")
	case r.Type == "":
		return fmt.Errorf("record type is required")
	case r.Name == "":
		return fmt.Errorf("record name is required")
	case r.Value == "":
		return fmt.Errorf("record value is required")
	}
	return nil
}

func (c *client) handleAPIResponse(body []byte, action string) error {
	var r Record
	err := json.Unmarshal(body, &r)
	if err != nil {
		return fmt.Errorf("unmarshal %s response: %w", action, err)
	}
	if r.Error != "" {
		return fmt.Errorf("API error during %s: %s", action, r.Error)
	}
	return nil
}
