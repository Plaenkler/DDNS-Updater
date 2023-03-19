package providers

import (
	"fmt"
)

type UpdateDNSPodRequest struct {
	Domain string
	Host   string
	Token  string
}

func UpdateDNSPod(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateDNSPodRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
