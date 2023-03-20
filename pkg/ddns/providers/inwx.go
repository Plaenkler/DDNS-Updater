package providers

import (
	"fmt"
)

type UpdateINWXRequest struct {
	Domain   string
	Host     string
	Username string
	Password string
}

func UpdateINWX(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateINWXRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
