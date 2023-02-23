package providers

import (
	"fmt"
)

type UpdateDynDNSRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateDynDNS(request interface{}) error {
	return fmt.Errorf("not implemented")
}
