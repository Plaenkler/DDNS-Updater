package providers

import (
	"fmt"
	"net/http"
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
	urlStr := fmt.Sprintf("https://%s:%s@my.maxihoster.com/nic/update?hostname=%s&myip=%s,%s", r.Username, r.Password, r.Host, ipAddr, "")
	_, err := SendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	return nil
}
