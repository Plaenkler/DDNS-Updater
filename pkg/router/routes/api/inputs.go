package api

import (
	"encoding/json"
	"net/http"

	"github.com/plaenkler/ddns/pkg/ddns"
)

func GetInputs(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		http.Error(w, "[api-GetInputs-1] missing provider", http.StatusBadRequest)
		return
	}
	request, ok := ddns.GetUpdaters()[provider]
	if !ok {
		http.Error(w, "[api-GetInputs-2] provider is not valid", http.StatusBadRequest)
		return
	}
	fields := &request.Request
	inputs, err := json.Marshal(fields)
	if err != nil {
		http.Error(w, "[api-GetInputs-3] could not marshal fields", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(inputs)
}
