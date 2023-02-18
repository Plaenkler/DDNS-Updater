package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/plaenkler/ddns/pkg/config"
	"github.com/plaenkler/ddns/pkg/model"
)

func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("[api-handleConfig-1] could not parse form err: %s", err.Error())
		}
		// 1 - 65535
		if len(r.FormValue("port")) < 1 || len(r.FormValue("port")) > 5 {
			log.Printf("[api-handleConfig-2] port is not valid")
		}
		if len(r.FormValue("interval")) < 1 {
			log.Printf("[api-handleConfig-3] interval is not valid")
		}
		port, err := strconv.ParseUint(r.FormValue("port"), 10, 16)
		if err != nil {
			log.Printf("[api-handleConfig-4] port is not valid - error: %s", err)
		}
		interval, err := strconv.ParseUint(r.FormValue("interval"), 10, 16)
		if err != nil {
			log.Printf("[api-handleConfig-5] interval is not valid - error: %s", err)
		}
		config.UpdateConfig(&model.Config{
			Port:     port,
			Interval: interval,
		})
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		if err != nil {
			log.Printf("[api-handleConfig-6] could not write http reply - error: %s", err)
		}
	}
}
