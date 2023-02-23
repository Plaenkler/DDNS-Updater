package providers

import (
	"fmt"
)

type UpdateCloudflareRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateCloudflare(request interface{}) error {
	return fmt.Errorf("not implemented")
}
