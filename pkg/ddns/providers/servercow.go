package providers

import (
	"fmt"
)

type UpdateServercowRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateServercow(request interface{}) error {
	return fmt.Errorf("not implemented")
}
