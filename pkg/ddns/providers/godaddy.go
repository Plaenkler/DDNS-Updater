package providers

import (
	"fmt"
)

type UpdateGoDaddyRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateGoDaddy(request interface{}) error {
	return fmt.Errorf("not implemented")
}
