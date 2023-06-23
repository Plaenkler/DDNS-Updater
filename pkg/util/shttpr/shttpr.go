package shttpr

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func SendHTTPRequest(method string, url string, auth *url.Userinfo) (string, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("[util-SendHTTPRequest-1] failed to create HTTP request: %v", err)
	}
	if auth != nil {
		password, _ := auth.Password()
		req.SetBasicAuth(auth.Username(), password)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("[util-SendHTTPRequest-2] failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("[util-SendHTTPRequest-3] HTTP request returned status code %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("[util-SendHTTPRequest-4] failed to read HTTP response body: %v", err)
	}
	return string(bytes.TrimSpace(body)), nil
}
