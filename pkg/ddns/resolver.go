package ddns

import (
	"fmt"
	"io"
	"net"
	"net/http"

	log "github.com/plaenkler/ddns-updater/pkg/logging"
)

var resolvers = map[string]string{
	"ipify":  "https://api.ipify.org",
	"my-ip":  "https://api.my-ip.io/ip",
	"ipych":  "https://api.ipy.ch",
	"intel":  "https://nms.intellitrend.de",
	"ident":  "https://ident.me/",
	"ifconf": "https://ifconfig.me/ip",
}

func GetPublicIP() (string, error) {
	for r := range resolvers {
		resp, err := http.Get(resolvers[r])
		if err != nil {
			log.Errorf("[ddns-GetPublicIP-1] %s failed: %s", r, err)
			continue
		}
		defer resp.Body.Close()
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("[ddns-GetPublicIP-2] %s failed: %s", r, err)
			continue
		}
		addr := string(bytes)
		if !isValidIPAddress(addr) {
			log.Errorf("[ddns-GetPublicIP-3] %s failed: %s", r, addr)
			continue
		}
		log.Infof("[ddns-GetPublicIP-4] %s succeeded: %s", r, addr)
		return addr, nil
	}
	return "", fmt.Errorf("[ddns-GetPublicIP-5] all resolvers failed")
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
