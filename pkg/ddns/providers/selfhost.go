package providers

import (
	"fmt"
)

type UpdateSelfhostRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateSelfhost(request interface{}) error {
	return fmt.Errorf("not implemented")
}
