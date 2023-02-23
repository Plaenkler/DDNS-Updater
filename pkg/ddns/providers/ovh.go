package providers

import (
	"fmt"
)

type UpdateOVHRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateOVH(request interface{}) error {
	return fmt.Errorf("not implemented")
}
