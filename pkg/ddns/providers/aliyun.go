package providers

import (
	"fmt"
)

type UpdateAliyunRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateAliyun(request interface{}) error {
	return fmt.Errorf("not implemented")
}
