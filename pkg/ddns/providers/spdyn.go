package providers

import (
	"fmt"
)

type UpdateSpdynRequest struct {
	Domain        string
	Host          string
	User          string
	Password      string
	Token         string
	UseProviderIP bool
}

func UpdateSpdyn(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateSpdynRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
