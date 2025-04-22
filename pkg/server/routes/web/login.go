package web

import (
	"fmt"
	"html/template"
	"net/http"

	log "github.com/plaenkler/ddns-updater/pkg/logging"
)

func ProvideLogin(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("login").ParseFS(static,
		"static/html/pages/login.html",
		"static/html/partials/include.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = fmt.Fprintf(w, "could not provide template: %s", err)
		if err != nil {
			log.Errorf("failed to write response: %v", err)
		}
		return
	}
	w.Header().Add("Content-Type", "text/html")
	err = tpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = fmt.Fprintf(w, "could not execute parsed template: %v", err)
		if err != nil {
			log.Errorf("failed to write response: %v", err)
		}
	}
}
