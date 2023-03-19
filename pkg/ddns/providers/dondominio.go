package providers

import (
	"fmt"
)

type UpdateDonDominioRequest struct {
	Domain   string
	Host     string
	Username string
	Password string
	Name     string
}

func UpdateDonDominio(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateDonDominioRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
