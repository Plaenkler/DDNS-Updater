package providers

import (
	"fmt"
)

type UpdateNoIPRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateNoIP(request interface{}) error {
	return fmt.Errorf("not implemented")
}
