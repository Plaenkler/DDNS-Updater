package providers

import (
	"fmt"
)

type UpdateNamecheapRequest struct {
	Domain        string
	Host          string
	Password      string
	UseProviderIP bool
}

func UpdateNamecheap(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateNamecheapRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
