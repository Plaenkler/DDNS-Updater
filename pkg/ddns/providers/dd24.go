package providers

import (
	"fmt"
	"net/http"
	"strings"
)

type UpdateDD24Request struct {
	Domain   string `json:"Domain"`
	Host     string `json:"Host"`
	Password string `json:"Password"`
}

func UpdateDD24(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateDD24Request)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	urlStr := fmt.Sprintf("https://dynamicdns.key-systems.net/update.php?hostname=%s&password=%s&ip=%s", r.Host, r.Password, ipAddr)
	resp, err := SendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	if strings.Contains(resp, "success") {
		return nil
	}
	return fmt.Errorf("failed to update DDNS entry: %s", resp)
}
