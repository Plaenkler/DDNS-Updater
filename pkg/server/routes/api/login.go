package api

import (
	"net/http"
	"time"

	"github.com/plaenkler/ddns-updater/pkg/config"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"github.com/plaenkler/ddns-updater/pkg/server/session"
	"github.com/plaenkler/ddns-updater/pkg/server/totps"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if !config.Get().UseTOTP {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	currentTOTP := r.FormValue("totp")
	if !totps.Verify(currentTOTP) {
		log.Errorf("invalid totp: %s", currentTOTP)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	token, err := session.Add()
	if err != nil {
		log.Errorf("could not add session: %s", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "ddns-updater-token",
		Value:    token,
		Expires:  time.Now().Add(10 * time.Minute),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
