package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/plaenkler/ddns-updater/pkg/cipher"
	"github.com/plaenkler/ddns-updater/pkg/database"
	"github.com/plaenkler/ddns-updater/pkg/database/model"
	"github.com/plaenkler/ddns-updater/pkg/ddns"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"gorm.io/gorm"
)

func CreateJob(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Errorf("[api-CreateJob-1] could not parse form: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	provider := r.FormValue("provider")
	updater, ok := ddns.GetUpdaters()[provider]
	if !ok {
		http.Error(w, "Invalid provider", http.StatusBadRequest)
		log.Errorf("[api-CreateJob-2] provider is not valid")
		return
	}
	jobModel := &updater.Request
	params := r.FormValue("params")
	err = json.Unmarshal([]byte(params), &jobModel)
	if err != nil {
		log.Errorf("[api-CreateJob-3] could not unmarshal params: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	encParams, err := cipher.Encrypt(params)
	if err != nil {
		log.Errorf("[api-CreateJob-4] could not encrypt params: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	job := model.SyncJob{
		Provider: provider,
		Params:   encParams,
	}
	db := database.GetDatabase()
	if db == nil {
		log.Errorf("[api-CreateJob-5] could not get database connection")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = db.Create(&job).Error
	if err != nil {
		log.Errorf("[api-CreateJob-6] could not create job: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Infof("[api-CreateJob-6] created job with ID %d", job.ID)
}

func UpdateJob(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		log.Errorf("[api-UpdateJob-1] could not parse form: %s", err)
		return
	}
	id, err := strconv.ParseUint(r.FormValue("ID"), 10, 32)
	if err != nil {
		http.Error(w, "ID is not valid", http.StatusBadRequest)
		log.Errorf("[api-UpdateJob-2] ID is not valid: %s", err)
		return
	}
	provider := r.FormValue("provider")
	updater, ok := ddns.GetUpdaters()[provider]
	if !ok {
		http.Error(w, "Invalid provider", http.StatusBadRequest)
		log.Errorf("[api-UpdateJob-3] provider is not valid")
		return
	}
	jobModel := &updater.Request
	params := r.FormValue("params")
	err = json.Unmarshal([]byte(params), &jobModel)
	if err != nil {
		http.Error(w, "Could not unmarshal params", http.StatusBadRequest)
		log.Errorf("[api-UpdateJob-4] could not unmarshal params: %s", err)
		return
	}
	encParams, err := cipher.Encrypt(params)
	if err != nil {
		http.Error(w, "Could not encrypt params", http.StatusInternalServerError)
		log.Errorf("[api-UpdateJob-5] could not encrypt params: %s", err)
		return
	}
	job := model.SyncJob{
		Model: gorm.Model{
			ID: uint(id),
		},
		Provider: provider,
		Params:   encParams,
	}
	db := database.GetDatabase()
	if db == nil {
		http.Error(w, "Could not get database connection", http.StatusInternalServerError)
		log.Errorf("[api-UpdateJob-6] could not get database connection")
		return
	}
	err = db.Save(&job).Error
	if err != nil {
		http.Error(w, "Could not update job", http.StatusInternalServerError)
		log.Errorf("[api-UpdateJob-7] could not update job: %s", err)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Infof("[api-UpdateJob-8] updated job with ID %d", job.ID)
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	strID := r.URL.Query().Get("ID")
	if len(strID) == 0 {
		log.Errorf("[api-DeleteJob-1] ID is not set")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(strID, 10, 32)
	if err != nil {
		log.Errorf("[api-DeleteJob-2] ID is not valid: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	job := model.SyncJob{
		Model: gorm.Model{
			ID: uint(id),
		},
	}
	db := database.GetDatabase()
	if db == nil {
		log.Errorf("[api-DeleteJob-3] could not get database connection")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := db.Unscoped().Delete(&job).Error; err != nil {
		log.Errorf("[api-DeleteJob-4] could not delete job: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Infof("[api-DeleteJob-5] deleted job with ID %d", job.ID)
}
