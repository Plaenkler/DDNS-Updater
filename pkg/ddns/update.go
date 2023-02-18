package ddns

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/plaenkler/ddns/pkg/model"
)

type Update func(job model.SyncJob, ipAddr string) error

var updaters = map[string]Update{
	"Strato": updateStrato,
	"DDNSS":  updateDDNSS,
}

func IsUpdaterSupported(updater string) bool {
	_, ok := updaters[updater]
	return ok
}

func sendHTTPRequest(method string, url string, auth *url.Userinfo) (string, error) {
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
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to update DDNS entry: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read HTTP response body: %v", err)
	}
	return string(bytes.TrimSpace(body)), nil
}

func updateStrato(job model.SyncJob, ipAddr string) error {
	urlStr := fmt.Sprintf("https://%s:%s@dyndns.strato.com/nic/update?hostname=%s&myip=%s", job.User, job.Password, job.Domain, ipAddr)
	resp, err := sendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	switch {
	case strings.HasPrefix(resp, "good "):
		return nil
	case strings.HasPrefix(resp, "nochg "):
		return nil
	default:
		return fmt.Errorf("failed to update DDNS entry: %s", resp)
	}
}

func updateDDNSS(job model.SyncJob, ipAddr string) error {
	urlStr := fmt.Sprintf("https://www.ddnss.de/upd.php?user=%s&pwd=%s&host=%s&ip=%s&ip6=%s", job.User, job.Password, job.Domain, ipAddr, "")
	resp, err := sendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	if strings.Contains(resp, "Error Occurred While Processing Request") {
		return fmt.Errorf("failed to update DDNS entry: %s", resp)
	}
	return nil
}
