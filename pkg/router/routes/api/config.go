package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/plaenkler/ddns/pkg/config"
)

func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("[api-handleConfig-1] could not parse form err: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(r.FormValue("interval")) < 1 {
		log.Printf("[api-handleConfig-2] interval is not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	port, err := strconv.ParseUint(r.FormValue("port"), 10, 16)
	if err != nil {
		log.Printf("[api-handleConfig-3] port is not valid - error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	interval, err := strconv.ParseUint(r.FormValue("interval"), 10, 16)
	if err != nil {
		log.Printf("[api-handleConfig-4] interval is not valid - error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	config.UpdateConfig(&config.Config{
		Port:     port,
		Interval: interval,
	})
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusTemporaryRedirect)
}
