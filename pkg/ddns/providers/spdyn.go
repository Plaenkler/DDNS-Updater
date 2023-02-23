package providers

import (
	"fmt"
)

type UpdateSpdynRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateSpdyn(request interface{}) error {
	return fmt.Errorf("not implemented")
}
