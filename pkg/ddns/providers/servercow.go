package providers

import (
	"fmt"
)

type UpdateServercowRequest struct {
	Username      string
	Host          string
	Domain        string
	Password      string
	UseProviderIP bool
	TTL           uint `json:",string"`
}

func UpdateServercow(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateServercowRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
