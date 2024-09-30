package web

import (
	"fmt"
	"html/template"
	"net/http"
)

func ProvideLogin(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("login").ParseFS(static,
		"static/html/pages/login.html",
		"static/html/partials/include.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not provide template: %s", err)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	err = tpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not execute parsed template: %v", err)
	}
}
