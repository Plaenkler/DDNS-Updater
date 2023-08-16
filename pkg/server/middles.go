package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/plaenkler/ddns-updater/pkg/server/session"
)

func limitRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := isOverLimit(r)
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

func controlCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		next.ServeHTTP(w, r)
	})
}

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedPrefixes := []string{
			"/css/",
			"/img/",
			"/login",
			"/api/login",
		}
		pth := r.URL.Path
		for _, path := range allowedPrefixes {
			if strings.HasPrefix(pth, path) {
				next.ServeHTTP(w, r)
				return
			}
		}
		c, err := r.Cookie("ddns-updater-token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if !session.Verify(c.Value) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
