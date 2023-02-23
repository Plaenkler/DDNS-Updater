package providers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/plaenkler/ddns/pkg/util"
)

type UpdateDDNSSRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
	IPAddr   string `json:"ipAddr"`
}

func UpdateDDNSS(request interface{}) error {
	r, ok := request.(UpdateDDNSSRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	urlStr := fmt.Sprintf("https://www.ddnss.de/upd.php?user=%s&pwd=%s&host=%s&ip=%s&ip6=%s", r.User, r.Password, r.Domain, r.IPAddr, "")
	resp, err := util.SendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	if strings.Contains(resp, "Error Occurred While Processing Request") {
		return fmt.Errorf("failed to update DDNS entry: %s", resp)
	}
	return nil
}
