package providers

import (
	"fmt"
)

type UpdateNjallaRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateNjalla(request interface{}) error {
	return fmt.Errorf("not implemented")
}
