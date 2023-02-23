package providers

import (
	"fmt"
)

type UpdateOpenDNSRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateOpenDNS(request interface{}) error {
	return fmt.Errorf("not implemented")
}
