package providers

import (
	"fmt"
)

type UpdateINWXRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateINWX(request interface{}) error {
	return fmt.Errorf("not implemented")
}
