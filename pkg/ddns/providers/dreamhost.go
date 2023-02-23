package providers

import (
	"fmt"
)

type UpdateDreamhostRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateDreamhost(request interface{}) error {
	return fmt.Errorf("not implemented")
}
