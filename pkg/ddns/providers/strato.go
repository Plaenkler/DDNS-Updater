package providers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/plaenkler/ddns/pkg/util"
)

type UpdateStratoRequest struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Domain   string `json:"Domain"`
}

func UpdateStrato(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateStratoRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	urlStr := fmt.Sprintf("https://%s:%s@dyndns.strato.com/nic/update?hostname=%s&myip=%s", r.User, r.Password, r.Domain, ipAddr)
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
