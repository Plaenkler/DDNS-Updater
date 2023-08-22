package providers

import (
	"encoding/json"
	"fmt"
)

type UpdateGCPRequest struct {
	Domain      string
	Host        string
	Project     string
	Zone        string
	Credentials json.RawMessage
}

func UpdateGCP(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateGCPRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Domain)
}
