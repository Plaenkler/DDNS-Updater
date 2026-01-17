package providers

import (
	"fmt"
	"net/http"
	"net/url"
)

type UpdateIONOSRequest struct {
	UpdateURL string
}

func UpdateIONOS(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateIONOSRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	if r.UpdateURL == "" {
		return fmt.Errorf("UpdateURL is required")
	}
	_, err := url.ParseRequestURI(r.UpdateURL)
	if err != nil {
		return fmt.Errorf("invalid UpdateURL: %w", err)
	}
	_, err = SendHTTPRequest(http.MethodGet, r.UpdateURL, nil)
	if err != nil {
		return err
	}
	return nil
}
