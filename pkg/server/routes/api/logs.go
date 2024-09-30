package api

import (
	"net/http"

	"github.com/plaenkler/ddns-updater/pkg/logging"
)

func GetLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := logging.GetEntries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
