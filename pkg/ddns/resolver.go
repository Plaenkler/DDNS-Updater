package ddns

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/plaenkler/ddns/pkg/model"
)

func sendHTTPRequest(method string, url string, auth *url.Userinfo) (*http.Response, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	if auth != nil {
		password, _ := auth.Password()
		req.SetBasicAuth(auth.Username(), password)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	return resp, nil
}

func updateStrato(job model.Updater, ipAddr string) error {
	u, _ := url.Parse(fmt.Sprintf("https://dyndns.strato.com/nic/update?system=dyndns&hostname=%s&myip=%s", job.Domain, ipAddr))
	resp, err := sendHTTPRequest(http.MethodGet, u.String(), url.UserPassword(job.User, job.Password))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP response body: %v", err)
	}
	response := string(bytes.TrimSpace(body))
	if !regexp.MustCompile("^good|^nochg").MatchString(response) {
		return fmt.Errorf("failed to update DDNS entry: %s", response)
	}
	return nil
}

func updateDDNSS(job model.Updater, ipAddr string) error {
	u, _ := url.Parse(fmt.Sprintf("https://www.ddnss.de/upd.php?key=%s&host=%s&ip=%s", job.Password, job.Domain, ipAddr))
	resp, err := sendHTTPRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP response body: %v", err)
	}
	response := string(bytes.TrimSpace(body))
	if !regexp.MustCompile("^good|^nochg").MatchString(response) {
		return fmt.Errorf("failed to update DDNS entry: %s", response)
	}
	return nil
}
