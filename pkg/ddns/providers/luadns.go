package providers

import (
	"fmt"
)

type UpdateLuaDNSRequest struct {
	Domain string
	Host   string
	EMail  string
	Token  string
}

func UpdateLuaDNS(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateLuaDNSRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
