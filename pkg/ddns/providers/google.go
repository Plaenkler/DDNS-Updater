package providers

import "fmt"

type UpdateGoogleRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateGoogle(request interface{}) error {
	return fmt.Errorf("not implemented")
}
