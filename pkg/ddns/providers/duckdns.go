package providers

import (
	"fmt"
)

type UpdateDuckDNSRequest struct {
	Host          string
	Token         string
	UseProviderIP bool
}

func UpdateDuckDNS(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateDuckDNSRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Host)
}
