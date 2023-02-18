package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/plaenkler/ddns/pkg/database"
	"github.com/plaenkler/ddns/pkg/ddns"
	"github.com/plaenkler/ddns/pkg/model"
	"golang.org/x/net/idna"
	"gorm.io/gorm"
)

func CreateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("[api-CreateJob-1] could not parse form - error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	provider := r.FormValue("provider")
	if !ddns.IsUpdaterSupported(provider) {
		log.Printf("[api-CreateJob-2] provider is not valid")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	domain := r.FormValue("domain")
	if _, err := idna.Lookup.ToASCII(domain); err != nil {
		log.Printf("[api-CreateJob-3] domain is not valid - error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user := r.FormValue("user")
	password := r.FormValue("password")
	if len(user) == 0 || len(password) == 0 {
		log.Printf("[api-CreateJob-4] user or password is not set")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	job := model.SyncJob{
		Provider: provider,
		Domain:   domain,
		User:     user,
		Password: password,
	}
	if err := database.GetManager().DB.Create(&job).Error; err != nil {
		log.Printf("[api-CreateJob-5] could not create job - error: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Printf("[api-CreateJob-6] created job with ID %d", job.ID)
}

func UpdateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		log.Printf("[api-UpdateJob-1] could not parse form - error: %s", err.Error())
		return
	}
	if !ddns.IsUpdaterSupported(r.FormValue("provider")) {
		http.Error(w, "Invalid provider", http.StatusBadRequest)
		log.Printf("[api-UpdateJob-2] provider is not valid")
		return
	}
	_, err = idna.Lookup.ToASCII(r.FormValue("domain"))
	if err != nil {
		http.Error(w, "Invalid domain", http.StatusBadRequest)
		log.Printf("[api-UpdateJob-3] domain is not valid - error: %s", err)
		return
	}
	if len(r.FormValue("user")) == 0 || len(r.FormValue("password")) == 0 {
		http.Error(w, "User or password not set", http.StatusBadRequest)
		log.Printf("[api-UpdateJob-4] user or password is not set")
		return
	}
	id, err := strconv.ParseUint(r.FormValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Printf("[api-UpdateJob-5] ID is not valid - error: %s", err)
		return
	}
	job := model.SyncJob{
		Model: gorm.Model{
			ID: uint(id),
		},
		Provider: r.FormValue("provider"),
		Domain:   r.FormValue("domain"),
		User:     r.FormValue("user"),
		Password: r.FormValue("password"),
	}
	err = database.GetManager().DB.Save(&job).Error
	if err != nil {
		http.Error(w, "Could not update job", http.StatusInternalServerError)
		log.Printf("[api-UpdateJob-6] could not update job - error: %s", err)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Printf("[api-UpdateJob-7] updated job with ID %d", job.ID)
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	strID := r.URL.Query().Get("ID")
	if len(strID) == 0 {
		log.Printf("[api-DeleteJob-1] ID is not set")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(strID, 10, 32)
	if err != nil {
		log.Printf("[api-DeleteJob-3] ID is not valid - error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	job := model.SyncJob{
		Model: gorm.Model{
			ID: uint(id),
		},
	}
	if err := database.GetManager().DB.Delete(&job).Error; err != nil {
		log.Printf("[api-DeleteJob-4] could not delete job - error: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Printf("[api-DeleteJob-5] deleted job with ID %d", job.ID)
}
