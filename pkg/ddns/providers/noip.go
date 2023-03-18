package providers

import (
	"fmt"
)

type UpdateNoIPRequest struct {
	Domain        string
	Host          string
	Username      string
	Password      string
	UseProviderIP bool
}

func UpdateNoIP(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateNoIPRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
