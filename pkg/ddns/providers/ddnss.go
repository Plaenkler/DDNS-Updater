package providers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/plaenkler/ddns/pkg/ddns/providers/shared"
)

type UpdateDDNSSRequest struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Domain   string `json:"Domain"`
}

func UpdateDDNSS(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateDDNSSRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	urlStr := fmt.Sprintf("https://www.ddnss.de/upd.php?user=%s&pwd=%s&host=%s&ip=%s&ip6=%s", r.User, r.Password, r.Domain, ipAddr, "")
	resp, err := shared.SendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	if strings.Contains(resp, "Error Occurred While Processing Request") {
		return fmt.Errorf("failed to update DDNS entry: %s", resp)
	}
	return nil
}
