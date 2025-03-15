package providers

import (
	"fmt"
	"net/http"
	"strings"
)

type UpdateFreeDNSRequest struct {
	Domain   string `json:"Domain"`
	Token string `json:"Token"`
}

func UpdateDD24(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateFreeDNSRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	urlStr := fmt.Sprintf("https://freedns.afraid.org/dynamic/update.php?hostname=%s&Token", r.Host, r.Token)
	resp, err := SendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	if strings.Contains(resp, "success") {
		return nil
	}
	return fmt.Errorf("failed to update DDNS entry: %s", resp)
}
