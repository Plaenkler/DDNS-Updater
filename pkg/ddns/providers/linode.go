package providers

import (
	"fmt"
)

type UpdateLinodeRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateLinode(request interface{}) error {
	return fmt.Errorf("not implemented")
}
