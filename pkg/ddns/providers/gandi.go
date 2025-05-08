package providers

import (
	"fmt"
)

type UpdateGandiRequest struct {
	Domain string
	Host   string
	TTL    int `json:",string"`
	Key    string
}

func UpdateGandi(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateGandiRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
