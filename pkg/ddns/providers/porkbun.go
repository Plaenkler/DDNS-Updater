package providers

import (
	"fmt"
)

type UpdatePorkbunRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdatePorkbun(request interface{}) error {
	return fmt.Errorf("not implemented")
}
