package providers

import (
	"fmt"
)

type UpdateVariomediaRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateVariomedia(request interface{}) error {
	return fmt.Errorf("not implemented")
}
