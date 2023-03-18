package providers

import (
	"fmt"
)

type UpdateLinodeRequest struct {
	Domain string
	Host   string
	Token  string
}

func UpdateLinode(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateLinodeRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
