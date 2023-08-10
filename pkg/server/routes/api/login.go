package api

import (
	"net/http"
	"time"

	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"github.com/plaenkler/ddns-updater/pkg/server/session"
	"github.com/plaenkler/ddns-updater/pkg/server/totp"
)

func Login(w http.ResponseWriter, r *http.Request) {
	currentTOTP := r.FormValue("totp")
	ok := totp.VerifiyTOTP(currentTOTP)
	if !ok {
		log.Errorf("[api-login-1] invalid totp: %s", currentTOTP)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	token, err := session.AddSession()
	if err != nil {
		log.Errorf("[api-login-2] could not add session: %s", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "ddns-updater-token",
		Value:    token,
		Expires:  time.Now().Add(10 * time.Minute),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
	http.Redirect(w, r, "/", http.StatusOK)
}
