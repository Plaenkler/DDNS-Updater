package providers

import (
	"fmt"
)

type UpdateDNSOMaticRequest struct {
	Domain        string
	Host          string
	Username      string
	Password      string
	UseProviderIP bool
}

func UpdateDNSOMatic(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateDNSOMaticRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
