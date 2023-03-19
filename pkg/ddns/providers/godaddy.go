package providers

import (
	"fmt"
)

type UpdateGoDaddyRequest struct {
	Domain string
	Host   string
	Key    string
	Secret string
}

func UpdateGoDaddy(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateGoDaddyRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
