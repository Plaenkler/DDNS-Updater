package providers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	log "github.com/plaenkler/ddns-updater/pkg/logging"
)

type UpdateInfomaniakRequest struct {
	Host     string
	Username string
	Password string
}

func UpdateInfomaniak(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateInfomaniakRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", r)
	}
	urlStr := fmt.Sprintf("https://%s:%s@infomaniak.com/nic/update?hostname=%s&myip=%s", r.Username, r.Password, r.Host, ipAddr)
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Plaenkler DDNS-Updater/V0 info@plaenkler.com")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Errorf("error closing response body: %v\n", err)
		}
	}()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	text := string(bytes)
	switch {
	case strings.Contains(text, "good"):
		return nil
	case strings.Contains(text, "nochg"):
		return nil
	default:
		return fmt.Errorf("failed to update DDNS entry: %s", text)
	}
}
