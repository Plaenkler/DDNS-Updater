package providers

import (
	"fmt"
)

type UpdateDNSOMaticRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateDNSOMatic(request interface{}) error {
	return fmt.Errorf("not implemented")
}
