package ddns

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var resolvers = map[string]string{
	"ipify": "https://api.ipify.org",
	"my-ip": "https://api.my-ip.io/ip",
}

func GetPublicIP() (string, error) {
	for r := range resolvers {
		resp, err := http.Get(resolvers[r])
		if err != nil {
			log.Printf("[resolver-1] %s failed: %s", r, err)
			continue
		}
		defer resp.Body.Close()
		ipBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[resolver-2] %s failed: %s", r, err)
			continue
		}
		return string(ipBytes), nil
	}
	return "", fmt.Errorf("[resolver-3] all resolvers failed")
}
