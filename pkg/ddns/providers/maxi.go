package providers

import (
	"fmt"
	"net/http"
	"strings"
)

type UpdateMaxiHosterRequest struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Host     string `json:"Host"`
}

func UpdateMaxiHoster(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateMaxiHosterRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	urlStr := fmt.Sprintf("https://my.maxihoster.com/dyndns.php?domain=%s&ip=%s&ipv6=%s&username=%s&password=%s", r.Host, ipAddr, "", r.Username, r.Password)
	resp, err := SendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	if strings.Contains(resp, "Error Occurred While Processing Request") {
		return fmt.Errorf("failed to update DDNS entry: %s", resp)
	}
	return nil
}
