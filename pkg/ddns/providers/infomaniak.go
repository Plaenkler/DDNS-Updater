package providers

import (
	"fmt"
)

type UpdateInfomaniakRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateInfomaniak(request interface{}) error {
	return fmt.Errorf("not implemented")
}
