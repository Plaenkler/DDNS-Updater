package providers

import (
	"fmt"
	"net/http"
	"strings"
)

type UpdateNamecheapRequest struct {
	Domain   string
	Host     string
	Password string
}

func UpdateNamecheap(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateNamecheapRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}

	urlStr := fmt.Sprintf("https://dynamicdns.park-your-domain.com/update?host=%s&domain=%s&password=%s&ip=%s", r.Host, r.Domain, r.Password, ipAddr)
	resp, err := SendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	if strings.Contains(resp, "success") {
		return nil
	}
	return fmt.Errorf("failed to update DDNS entry: %s", resp)
}
