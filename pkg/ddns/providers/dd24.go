package providers

import (
	"fmt"
)

type UpdateDD24Request struct {
	Domain        string
	Host          string
	Password      string
	UseProviderIP bool
}

func UpdateDD24(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateDD24Request)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
