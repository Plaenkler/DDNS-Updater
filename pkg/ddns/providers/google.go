package providers

import "fmt"

type UpdateGoogleRequest struct {
	Domain        string
	host          string
	Username      string
	Password      string
	UseProviderIP bool
}

func UpdateGoogle(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateGoogleRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
