package providers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	log "github.com/plaenkler/ddns-updater/pkg/logging"
)

func SendHTTPRequest(method string, url string, auth *url.Userinfo) (string, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}
	if auth != nil {
		password, _ := auth.Password()
		req.SetBasicAuth(auth.Username(), password)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer log.ErrorClose(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request returned status code %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read HTTP response body: %v", err)
	}
	return string(bytes.TrimSpace(body)), nil
}
