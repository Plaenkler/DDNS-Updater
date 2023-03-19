package providers

import (
	"fmt"
)

type UpdateSelfhostRequest struct {
	Domain        string
	Gost          string
	Username      string
	Password      string
	UseProviderIP bool
}

func UpdateSelfhost(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateSelfhostRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
