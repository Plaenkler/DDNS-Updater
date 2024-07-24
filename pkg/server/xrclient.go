package server

import (
	"net"
	"net/http"
	"strings"
)

const (
	xForwardedFor  = "X-Forwarded-For"
	xRealIP        = "X-Real-IP"
	xForwardedPort = "X-Forwarded-Port"
)

func getRealClientIP(r *http.Request) (string, error) {
	rAddr := r.Header.Get(xRealIP)
	if rAddr != "" {
		return rAddr, nil
	}
	rAddr = crawlForwardedIPs(r)
	if rAddr != "" {
		return rAddr, nil
	}
	rAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
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
