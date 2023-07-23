package ddns

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/plaenkler/ddns-updater/pkg/config"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
)

var resolvers = map[string]string{
	"ipify":  "https://api.ipify.org",
	"my-ip":  "https://api.my-ip.io/ip",
	"ipych":  "https://api.ipy.ch",
	"intel":  "https://nms.intellitrend.de",
	"ident":  "https://ident.me/",
	"ifconf": "https://ifconfig.me/ip",
	"icanh":  "https://icanhazip.com/",
}

func GetPublicIP() (string, error) {
	cRes := config.GetConfig().Resolver
	if cRes != "" {
		addr, err := resolveIPAddress(cRes)
		if err != nil {
			return "", fmt.Errorf("[ddns-GetPublicIP-1] resolver %s failed: %s", cRes, err)
		}
		log.Infof("[ddns-GetPublicIP-2] %s succeeded: %s", cRes, addr)
		return addr, nil
	}
	for r := range resolvers {
		addr, err := resolveIPAddress(resolvers[r])
		if err != nil {
			log.Errorf("[ddns-GetPublicIP-3] resolver %s failed: %s", r, err)
			continue
		}
		log.Infof("[ddns-GetPublicIP-4] %s succeeded: %s", r, addr)
		return addr, nil
	}
	return "", fmt.Errorf("[ddns-GetPublicIP-5] all resolvers failed")
}

func resolveIPAddress(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	addr := string(bytes)
	if !isValidIPAddress(addr) {
		return "", err
	}
	return addr, nil
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
