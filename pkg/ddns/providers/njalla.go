package providers

import (
	"fmt"
)

type UpdateNjallaRequest struct {
	Domain        string
	Host          string
	Key           string
	UseProviderIP bool
}

func UpdateNjalla(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateNjallaRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
