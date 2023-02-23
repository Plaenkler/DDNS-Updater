package providers

import (
	"fmt"
)

type UpdateDD24Request struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateDD24(request interface{}) error {
	return fmt.Errorf("not implemented")
}
