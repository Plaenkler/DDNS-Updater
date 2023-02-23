package providers

import (
	"fmt"
)

type UpdateAllInklRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateAllInkl(request interface{}) error {
	return fmt.Errorf("not implemented")
}
