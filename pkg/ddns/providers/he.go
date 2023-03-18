package providers

import (
	"fmt"
)

type UpdateHeRequest struct {
	Domain        string
	Host          string
	Password      string
	UseProviderIP bool
}

func UpdateHe(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateHeRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
