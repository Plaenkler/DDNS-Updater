package providers

import (
	"fmt"
	"net/http"
	"strings"
)

type UpdateCustomRequest struct {
	Domain string `json:"URL"`
	Check  string `json:"Check"`
}

func UpdateCustom(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateCustomRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	if !strings.Contains(r.Domain, "<ipv4>") {
		return fmt.Errorf("no <ipv4> placeholder found in URL")
	}
	r.Domain = strings.Replace(r.Domain, "<ipv4>", ipAddr, -1)
	resp, err := SendHTTPRequest(http.MethodGet, r.Domain, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	if !strings.Contains(string(body), r.Check) {
		return fmt.Errorf("check string '%s' not found in response", r.Check)
	}
	return nil
}
