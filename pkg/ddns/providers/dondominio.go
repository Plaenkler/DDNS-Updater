package providers

import (
	"fmt"
)

type UpdateDonDominioRequest struct {
	IPAddr string `json:"ipAddr"`
}

func UpdateDonDominio(request interface{}) error {
	return fmt.Errorf("not implemented")
}
