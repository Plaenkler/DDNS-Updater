package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/plaenkler/ddns/pkg/database"
	"github.com/plaenkler/ddns/pkg/database/model"
	"github.com/plaenkler/ddns/pkg/ddns"
	"github.com/plaenkler/ddns/pkg/util/limit"
	"gorm.io/gorm"
)

func CreateJob(w http.ResponseWriter, r *http.Request) {
	err := limit.IsOverLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Printf("[api-CreateJob-1] could not parse form - error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	provider := r.FormValue("provider")
	updater, ok := ddns.GetUpdaters()[provider]
	if !ok {
		http.Error(w, "Invalid provider", http.StatusBadRequest)
		log.Printf("[api-CreateJob-2] provider is not valid")
		return
	}
	jobModel := &updater.Request
	params := r.FormValue("params")
	err = json.Unmarshal([]byte(params), &jobModel)
	if err != nil {
		log.Printf("[api-CreateJob-3] could not unmarshal params - error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	job := model.SyncJob{
		Provider: provider,
		Params:   params,
	}
	err = database.GetManager().DB.Create(&job).Error
	if err != nil {
		log.Printf("[api-CreateJob-4] could not create job - error: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Printf("[api-CreateJob-5] created job with ID %d", job.ID)
}

func UpdateJob(w http.ResponseWriter, r *http.Request) {
	err := limit.IsOverLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		log.Printf("[api-UpdateJob-1] could not parse form - error: %s", err)
		return
	}
	id, err := strconv.ParseUint(r.FormValue("ID"), 10, 32)
	if err != nil {
		http.Error(w, "ID is not valid", http.StatusBadRequest)
		log.Printf("[api-UpdateJob-2] ID is not valid - error: %s", err)
		return
	}
	provider := r.FormValue("provider")
	updater, ok := ddns.GetUpdaters()[provider]
	if !ok {
		http.Error(w, "Invalid provider", http.StatusBadRequest)
		log.Printf("[api-UpdateJob-3] provider is not valid")
		return
	}
	jobModel := &updater.Request
	params := r.FormValue("params")
	err = json.Unmarshal([]byte(params), &jobModel)
	if err != nil {
		http.Error(w, "Could not unmarshal params", http.StatusBadRequest)
		log.Printf("[api-UpdateJob-4] could not unmarshal params - error: %s", err)
		return
	}
	job := model.SyncJob{
		Model: gorm.Model{
			ID: uint(id),
		},
		Provider: provider,
		Params:   params,
	}
	err = database.GetManager().DB.Save(&job).Error
	if err != nil {
		http.Error(w, "Could not update job", http.StatusInternalServerError)
		log.Printf("[api-UpdateJob-5] could not update job - error: %s", err)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Printf("[api-UpdateJob-6] updated job with ID %d", job.ID)
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	err := limit.IsOverLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	strID := r.URL.Query().Get("ID")
	if len(strID) == 0 {
		log.Printf("[api-DeleteJob-1] ID is not set")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(strID, 10, 32)
	if err != nil {
		log.Printf("[api-DeleteJob-2] ID is not valid - error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	job := model.SyncJob{
		Model: gorm.Model{
			ID: uint(id),
		},
	}
	if err := database.GetManager().DB.Unscoped().Delete(&job).Error; err != nil {
		log.Printf("[api-DeleteJob-3] could not delete job - error: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Printf("[api-DeleteJob-4] deleted job with ID %d", job.ID)
}
