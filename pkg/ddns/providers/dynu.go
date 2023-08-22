package providers

import (
	"fmt"
	"net/http"
	"strings"
)

type UpdateDynuRequest struct {
	User     string `json:"User"`
	Password string `json:"Password"`
}

func UpdateDynu(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateDynuRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	urlStr := fmt.Sprintf("https://%s:%s@api.dynu.com/nic/update?myip=%s&myipv6=", r.User, r.Password, ipAddr)
	resp, err := SendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	switch {
	case strings.HasPrefix(resp, "good"):
		return nil
	case strings.HasPrefix(resp, "nochg"):
		return nil
	default:
		return fmt.Errorf("failed to update DDNS entry: %s", resp)
	}
}
