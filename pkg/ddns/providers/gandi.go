package providers

import (
	"fmt"
)

type UpdateGandiRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateGandi(request interface{}) error {
	return fmt.Errorf("not implemented")
}
