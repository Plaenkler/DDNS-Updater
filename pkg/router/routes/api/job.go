package api

import (
	"log"
	"net/http"

	"github.com/plaenkler/ddns/pkg/database"
	"github.com/plaenkler/ddns/pkg/ddns"
	"github.com/plaenkler/ddns/pkg/model"
	"golang.org/x/net/idna"
)

func CreateJob(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("[api-CreateJob-1] could not parse form - error: %s", err.Error())
		}
		if !ddns.IsUpdaterSupported(r.FormValue("provider")) {
			log.Printf("[api-CreateJob-2] provider is not valid")
		}
		_, err = idna.Lookup.ToASCII(r.FormValue("domain"))
		if err != nil {
			log.Printf("[api-CreateJob-3] domain is not valid - error: %s", err)
		}
		if len(r.FormValue("user")) == 0 || len(r.FormValue("password")) == 0 {
			log.Printf("[api-CreateJob-4] user or password is not set")
		}
		job := model.SyncJob{
			Provider: r.FormValue("provider"),
			Domain:   r.FormValue("domain"),
			User:     r.FormValue("user"),
			Password: r.FormValue("password"),
		}
		err = database.GetManager().DB.FirstOrCreate(&job).Error
		if err != nil {
			log.Printf("[api-CreateJob-5] could not create job - error: %s", err)
		}
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		if err != nil {
			log.Printf("[api-CreateJob-6] could not write http reply - error: %s", err)
		}
	}
}

func UpdateJob(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Parse the request body into a new job
		// Update the existing job
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		if err != nil {
			log.Printf("[api-UpdateJob-1] could not write http reply - error: %s", err)
		}
	}
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Delete the job
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		if err != nil {
			log.Printf("[api-DeleteJob-1] could not write http reply - error: %s", err)
		}
	}
}
