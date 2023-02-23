package providers

import (
	"fmt"
)

type UpdateLuaDNSRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateLuaDNS(request interface{}) error {
	return fmt.Errorf("not implemented")
}
