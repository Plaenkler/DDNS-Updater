package providers

import (
	"fmt"
	"net/url"
)

type UpdateOVHRequest struct {
	Domain        string
	Host          string
	Username      string
	Password      string
	UseProviderIP bool
	Mode          string
	APIURL        *url.URL
	AppKey        string
	AppSecret     string
	ConsumerKey   string
}

func UpdateOVH(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateOVHRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
