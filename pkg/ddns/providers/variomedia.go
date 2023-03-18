package providers

import (
	"fmt"
)

type UpdateVariomediaRequest struct {
	Domain        string
	Host          string
	EMail         string
	Password      string
	UseProviderIP bool
}

func UpdateVariomedia(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateVariomediaRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
