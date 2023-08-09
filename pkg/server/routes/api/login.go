package api

import (
	"net/http"

	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"github.com/plaenkler/ddns-updater/pkg/server/session"
	"github.com/plaenkler/ddns-updater/pkg/server/totp"
)

func Login(w http.ResponseWriter, r *http.Request) {
	currentTOTP := r.FormValue("totp")
	ok := totp.VerifiyTOTP(currentTOTP)
	if !ok {
		log.Errorf("[api-login-1] invalid totp: %s", currentTOTP)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	token, err := session.AddSession()
	if err != nil {
		log.Errorf("[api-login-2] could not add session: %s", err)
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "ddns-updater-token",
		Value: token,
		Path:  "/",
	})
	http.Redirect(w, r, "/", http.StatusOK)
}
