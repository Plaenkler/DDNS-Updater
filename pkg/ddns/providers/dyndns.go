package providers

import (
	"fmt"
)

type UpdateDynDNSRequest struct {
	Domain        string
	Host          string
	Username      string
	ClientKey     string
	UseProviderIP bool
}

func UpdateDynDNS(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateDynDNSRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
