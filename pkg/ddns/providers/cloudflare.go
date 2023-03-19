package providers

import (
	"fmt"
)

type UpdateCloudflareRequest struct {
	Domain         string
	Host           string
	Key            string
	Token          string
	EMail          string
	UserServiceKey string
	ZoneIdentifier string
	Proxied        bool
	TTL            uint
}

func UpdateCloudflare(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateCloudflareRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
