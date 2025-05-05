package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	hetznerBaseURL    = "https://dns.hetzner.com/api/v1"
	updateRecordURI   = "/records/%s"
	headerContentType = "Content-Type"
	headerAPIToken    = "Auth-API-Token"
	contentTypeJSON   = "application/json"
)

type record struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	TTL    int    `json:"ttl"`
	ZoneID string `json:"zone_id"`
}

type client struct {
	APIToken string
}

type UpdateHetznerRequest struct {
	APIToken     string
	RecordID     string
	RecordType   string
	RecordName   string
	RecordValue  string
	RecordTTL    int
	RecordZoneID string
}

func (c *client) put(url string, body []byte) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	req.Header.Set(headerContentType, contentTypeJSON)
	req.Header.Set(headerAPIToken, c.APIToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *client) updateRecord(record record) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(hetznerBaseURL+updateRecordURI, record.ID)
	_, err = c.put(url, data)
	return err
}

func UpdateHetzner(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateHetznerRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	rec := record{
		ID:     r.RecordID,
		Type:   r.RecordType,
		Name:   r.RecordName,
		Value:  ipAddr,
		TTL:    r.RecordTTL,
		ZoneID: r.RecordZoneID,
	}
	c := &client{
		APIToken: r.APIToken,
	}
	err := c.updateRecord(rec)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}
	return nil
}
