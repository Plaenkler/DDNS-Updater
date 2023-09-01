package providers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type UpdateINWXRequest struct {
	Domain   string `json:"Domain"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func UpdateINWX(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateINWXRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	urlStr := fmt.Sprintf("https://dyndns.inwx.com/nic/update?hostname=%s&myip=%s", r.Domain, ipAddr)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(r.Username, r.Password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respStr := string(bodyBytes)
	switch {
	case strings.Contains(respStr, "good"):
		return nil
	case strings.Contains(respStr, "nochg"):
		return nil
	default:
		return fmt.Errorf("failed to update DDNS entry: %s", respStr)
	}
}
