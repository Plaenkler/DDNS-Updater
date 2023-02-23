package providers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/plaenkler/ddns/pkg/util"
)

type UpdateStratoRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
	IPAddr   string `json:"ipAddr"`
}

func UpdateStrato(request interface{}) error {
	r, ok := request.(UpdateStratoRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	urlStr := fmt.Sprintf("https://%s:%s@dyndns.strato.com/nic/update?hostname=%s&myip=%s", r.User, r.Password, r.Domain, r.IPAddr)
	resp, err := util.SendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	switch {
	case strings.Contains(resp, "good"):
		return nil
	case strings.Contains(resp, "nochg"):
		return nil
	default:
		return fmt.Errorf("failed to update DDNS entry: %s", resp)
	}
}
