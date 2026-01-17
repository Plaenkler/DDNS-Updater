package providers

import (
	"fmt"
	"net/http"
)

type UpdateIONOSRequest struct {
	UpdateURL string
}

func UpdateIONOS(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateIONOSRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	_, err := SendHTTPRequest(http.MethodGet, r.UpdateURL, nil)
	if err != nil {
		return err
	}
	return nil
}
