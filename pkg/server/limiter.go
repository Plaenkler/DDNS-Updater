package server

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter *rate.Limiter
	address string
}

var ipLimiters = map[string]ipLimiter{}

func IsOverLimit(r *http.Request) error {
	addr, err := getRealClientIP(r)
	if err != nil {
		return fmt.Errorf("[server-IsOverLimit-1] could not get client ip address")
	}
	iplm, ok := ipLimiters[addr]
	if !ok {
		iplm = ipLimiter{
			limiter: rate.NewLimiter(1, 3),
			address: addr,
		}
		ipLimiters[addr] = iplm
	}
	if !iplm.limiter.Allow() {
		return fmt.Errorf("[server-IsOverLimit-2] ip address %s is over limit", addr)
	}
	return nil
}

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
