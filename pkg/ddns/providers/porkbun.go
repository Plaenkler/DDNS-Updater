package providers

import (
	"fmt"
)

type UpdatePorkbunRequest struct {
	Domain       string
	Host         string
	TTL          uint `json:",string"`
	APIKey       string
	SecretAPIKey string
}

func UpdatePorkbun(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdatePorkbunRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
