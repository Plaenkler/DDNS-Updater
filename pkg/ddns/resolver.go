package ddns

import (
	"io"
	"net/http"
)

func GetPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ipBytes), nil
}