package providers

import (
	"fmt"
)

type UpdateInfomaniakRequest struct {
	Domain        string
	Host          string
	Username      string
	Password      string
	UseProviderIP bool
}

func UpdateInfomaniak(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateInfomaniakRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
