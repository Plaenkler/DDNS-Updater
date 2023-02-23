package providers

import (
	"fmt"
)

type UpdateDNSPodRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateDNSPod(request interface{}) error {
	return fmt.Errorf("not implemented")
}
