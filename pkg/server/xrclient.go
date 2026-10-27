package server

import (
	"net"
	"net/http"
	"strings"

	"github.com/plaenkler/ddns-updater/pkg/config"
)

const (
	xForwardedFor  = "X-Forwarded-For"
	xRealIP        = "X-Real-IP"
	xForwardedPort = "X-Forwarded-Port"
)

func getRealClientIP(r *http.Request) (string, error) {
	rAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	if !config.Get().TrustProxy {
		return rAddr, nil
	}
	pAddr := r.Header.Get(xRealIP)
	if pAddr != "" {
		return pAddr, nil
	}
	pAddr = crawlForwardedIPs(r)
	if pAddr != "" {
		return pAddr, nil
	}
	return rAddr, nil
}

func crawlForwardedIPs(r *http.Request) string {
	xff := r.Header.Get(xForwardedFor)
	ips := strings.Split(xff, ",")
	// Check IP addresses in reverse order
	for i := len(ips) - 1; i >= 0; i-- {
		ip := strings.TrimSpace(ips[i])
		if ip != "" {
			return ip
		}
	}
	return ""
}
