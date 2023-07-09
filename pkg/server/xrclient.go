package server

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func getRealClientIP(r *http.Request) (string, error) {
	clientIP := r.Header.Get("X-Real-IP")
	if clientIP != "" {
		return clientIP, nil
	}
	xff := r.Header.Get("X-Forwarded-For")
	ips := strings.Split(xff, ",")
	for i := len(ips) - 1; i >= 0; i-- {
		// Check IP addresses in reverse order to find real IP
		ip := strings.TrimSpace(ips[i])
		if ip != "" {
			clientIP = ip
			break
		}
	}
	if clientIP != "" {
		return clientIP, nil
	}
	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return clientIP, nil
}

func getRealClientPort(r *http.Request) (string, error) {
	remotePort := r.Header.Get("X-Forwarded-Port")
	if remotePort != "" {
		return remotePort, nil
	}
	// Extract remote address from TCP connection
	remoteAddr := r.RemoteAddr
	if !strings.Contains(remoteAddr, ":") {
		return "", fmt.Errorf("[server-GetRealClientPort-2] could not determine remote port")
	}
	_, port, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return "", err
	}
	return port, nil
}
