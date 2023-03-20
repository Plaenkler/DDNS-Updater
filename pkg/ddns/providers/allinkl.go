package providers

import (
	"fmt"
)

type UpdateAllInklRequest struct {
	Domain        string
	Host          string
	Username      string
	Password      string
	UseProviderIP bool
}

func UpdateAllInkl(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateAllInklRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
