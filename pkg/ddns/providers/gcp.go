package providers

import (
	"fmt"
)

type UpdateGCPRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateGCP(request interface{}) error {
	return fmt.Errorf("not implemented")
}
