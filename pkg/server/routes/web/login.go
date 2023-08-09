package web

import (
	"fmt"
	"html/template"
	"net/http"
)

func ProvideLogin(w http.ResponseWriter, r *http.Request) {
	template, err := template.New("login").ParseFS(static,
		"static/html/pages/login.html",
		"static/html/partials/include.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideLogin-1] could not provide template - error: %s", err)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	err = template.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideLogin-2] could not execute parsed template - error: %v", err)
	}
}
