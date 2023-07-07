package server

import (
	"fmt"
	"net/http"
)

func limitRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := IsOverLimit(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func forwardToProxy(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr, err := getRealClientIP(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		port, err := getRealClientPort(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.URL.Host = fmt.Sprintf("%s:%s", addr, port)
		handler.ServeHTTP(w, r)
	})
}
