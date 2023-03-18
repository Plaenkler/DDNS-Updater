package providers

import (
	"fmt"
)

type UpdateDigitalOceanRequest struct {
	Domain string
	Host   string
	Token  string
}

func UpdateDigitalOcean(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateDigitalOceanRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
