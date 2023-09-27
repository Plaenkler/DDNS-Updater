package server

import (
	"fmt"
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

func getRealClientPort(r *http.Request) (string, error) {
	remotePort := r.Header.Get(xForwardedPort)
	if remotePort != "" {
		return remotePort, nil
	}
	if !strings.Contains(r.RemoteAddr, ":") {
		return "", fmt.Errorf("[server-GetRealClientPort-2] could not determine remote port")
	}
	_, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return port, nil
}
