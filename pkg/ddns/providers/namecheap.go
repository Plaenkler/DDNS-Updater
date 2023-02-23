package providers

import (
	"fmt"
)

type UpdateNamecheapRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateNamecheap(request interface{}) error {
	return fmt.Errorf("not implemented")
}
