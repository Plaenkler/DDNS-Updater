package providers

import (
	"fmt"
)

type UpdateOpenDNSRequest struct {
	Domain        string
	Host          string
	Username      string
	Password      string
	UseProviderIP bool
}

func UpdateOpenDNS(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateOpenDNSRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
