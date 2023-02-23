package api

import (
	"net/http"

	"github.com/plaenkler/ddns/pkg/ddns"
	"github.com/plaenkler/ddns/pkg/util"
)

func GetInputs(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	request, ok := ddns.GetUpdaters()[provider]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(util.StructToHTML(request.Request)))
}
