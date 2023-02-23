package providers

import (
	"fmt"
)

type UpdateDigitalOceanRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateDigitalOcean(request interface{}) error {
	return fmt.Errorf("not implemented")
}
