package providers

import (
	"fmt"
)

type UpdateDreamhostRequest struct {
	Domain string
	Host   string
	Key    string
}

func UpdateDreamhost(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateDreamhostRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
