package providers

import (
	"fmt"
)

type UpdateDuckDNSRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateDuckDNS(request interface{}) error {
	return fmt.Errorf("not implemented")
}
