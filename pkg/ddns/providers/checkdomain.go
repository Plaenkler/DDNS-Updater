package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/plaenkler/ddns-updater/pkg/logging"
)

const (
	checkdomainBaseURL = "https://api.checkdomain.de/v1/dns/records"
	checkdomainTimeout = 10 * time.Second
)

type UpdateCheckdomainRequest struct {
	APIToken   string `json:"APIToken"`
	Domain     string `json:"Domain"`
	RecordName string `json:"RecordName"`
	RecordType string `json:"RecordType"`
	TTL        int    `json:"TTL,string"`
}

type checkdomainRecord struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	TTL    int    `json:"ttl,omitempty"`
	Domain string `json:"domain,omitempty"`
}

type checkdomainResponse struct {
	Records []checkdomainRecord `json:"records,omitempty"`
	Error   string              `json:"error,omitempty"`
	Message string              `json:"message,omitempty"`
}

func UpdateCheckdomain(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateCheckdomainRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}

	client := &http.Client{
		Timeout: checkdomainTimeout,
	}

	// Set default values
	if r.RecordType == "" {
		r.RecordType = "A"
	}
	if r.TTL == 0 {
		r.TTL = 3600
	}

	// Find existing record
	records, err := getCheckdomainRecords(client, r.APIToken, r.Domain, r.RecordName, r.RecordType)
	if err != nil {
		return fmt.Errorf("failed to get records: %w", err)
	}

	record := checkdomainRecord{
		Name:   r.RecordName,
		Type:   r.RecordType,
		Value:  ipAddr,
		TTL:    r.TTL,
		Domain: r.Domain,
	}

	if len(records) > 0 {
		// Update existing record
		record.ID = records[0].ID
		return updateCheckdomainRecord(client, r.APIToken, &record)
	}

	// Create new record
	return createCheckdomainRecord(client, r.APIToken, &record)
}

func getCheckdomainRecords(client *http.Client, token, domain, name, recordType string) ([]checkdomainRecord, error) {
	url := fmt.Sprintf("%s?domain=%s&name=%s&type=%s", checkdomainBaseURL, domain, name, recordType)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer log.ErrorClose(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response checkdomainResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("API error: %s", response.Error)
	}

	return response.Records, nil
}

func createCheckdomainRecord(client *http.Client, token string, record *checkdomainRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, checkdomainBaseURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer log.ErrorClose(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create record, status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func updateCheckdomainRecord(client *http.Client, token string, record *checkdomainRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	url := fmt.Sprintf("%s/%s", checkdomainBaseURL, record.ID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer log.ErrorClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update record, status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
