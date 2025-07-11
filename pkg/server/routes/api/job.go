package api

import (
	"encoding/json"
	"fmt"
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
		log.Errorf("could not parse form: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	provider := r.FormValue("provider")
	params := r.FormValue("params")
	err = verifyJobModel(provider, params)
	if err != nil {
		http.Error(w, "Invalid values", http.StatusBadRequest)
		log.Errorf("invalid values: %s", err)
		return
	}
	encParams, err := cipher.Encrypt(params)
	if err != nil {
		log.Errorf("could not encrypt params: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	job := model.SyncJob{
		Provider: provider,
		Params:   encParams,
	}
	db := database.GetDatabase()
	if db == nil {
		log.Errorf("could not get database connection")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = db.Create(&job).Error
	if err != nil {
		log.Errorf("could not create job: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Infof("created job with ID %d", job.ID)
}

func UpdateJob(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		log.Errorf("could not parse form: %s", err)
		return
	}
	id, err := strconv.ParseUint(r.FormValue("ID"), 10, 32)
	if err != nil {
		http.Error(w, "ID is not valid", http.StatusBadRequest)
		log.Errorf("ID is not valid: %s", err)
		return
	}
	provider := r.FormValue("provider")
	params := r.FormValue("params")
	err = verifyJobModel(provider, params)
	if err != nil {
		http.Error(w, "Invalid values", http.StatusBadRequest)
		log.Errorf("invalid values: %s", err)
		return
	}
	encParams, err := cipher.Encrypt(params)
	if err != nil {
		http.Error(w, "Could not encrypt params", http.StatusInternalServerError)
		log.Errorf("could not encrypt params: %s", err)
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
		log.Errorf("could not get database connection")
		return
	}
	err = db.Save(&job).Error
	if err != nil {
		http.Error(w, "Could not update job", http.StatusInternalServerError)
		log.Errorf("could not update job: %s", err)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Infof("updated job with ID %d", job.ID)
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	strID := r.URL.Query().Get("ID")
	if len(strID) == 0 {
		log.Errorf("ID is not set")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(strID, 10, 32)
	if err != nil {
		log.Errorf("ID is not valid: %s", err)
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
		log.Errorf("could not get database connection")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := db.Unscoped().Delete(&job).Error; err != nil {
		log.Errorf("could not delete job: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	log.Infof("deleted job with ID %d", job.ID)
}

func verifyJobModel(provider, params string) error {
	updater, ok := ddns.GetUpdaters()[provider]
	if !ok {
		return fmt.Errorf("provider is not valid")
	}
	jobModel := updater.Factory()
	return json.Unmarshal([]byte(params), jobModel)
}
