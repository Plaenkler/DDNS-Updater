package ddns

import (
	"fmt"
	"io"
	"log"
	"net"
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
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[resolver-2] %s failed: %s", r, err)
			continue
		}
		addr := string(bytes)
		if !isValidIPAddress(addr) {
			log.Printf("[resolver-3] %s failed: %s", r, addr)
			continue
		}
		log.Printf("[resolver-4] %s succeeded: %s", r, addr)
		return addr, nil
	}
	return "", fmt.Errorf("[resolver-5] all resolvers failed")
}

func isValidIPAddress(ip string) bool {
	addr := net.ParseIP(ip)
	if addr == nil {
		return false
	}
	if addr.IsUnspecified() {
		return false
	}
	if addr.IsPrivate() {
		return false
	}
	if addr.IsLoopback() {
		return false
	}
	if addr.IsMulticast() {
		return false
	}
	return true
}
