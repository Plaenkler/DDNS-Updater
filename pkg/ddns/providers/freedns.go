package providers

import (
	"fmt"
)

type UpdateFreeDNSRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateFreeDNS(request interface{}) error {
	return fmt.Errorf("not implemented")
}
