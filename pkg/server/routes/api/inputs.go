package api

import (
	"encoding/json"
	"net/http"

	"github.com/plaenkler/ddns-updater/pkg/ddns"
)

func GetInputs(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		http.Error(w, "missing provider", http.StatusBadRequest)
		return
	}
	updater, ok := ddns.GetUpdaters()[provider]
	if !ok {
		http.Error(w, "provider is not valid", http.StatusBadRequest)
		return
	}
	fields := updater.Factory()
	inputs, err := json.Marshal(fields)
	if err != nil {
		http.Error(w, "could not marshal fields", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(inputs)
	if err != nil {
		http.Error(w, "could not write response", http.StatusInternalServerError)
		return
	}
}
