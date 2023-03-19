package providers

import (
	"fmt"
)

type UpdateFreeDNSRequest struct {
	Domain string
	Host   string
	Token  string
}

func UpdateFreeDNS(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateFreeDNSRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
