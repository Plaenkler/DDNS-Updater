package providers

import (
	"fmt"
)

type UpdateHeRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateHe(request interface{}) error {
	return fmt.Errorf("not implemented")
}
