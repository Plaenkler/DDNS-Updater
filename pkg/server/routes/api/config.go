package api

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/plaenkler/ddns-updater/pkg/config"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"github.com/plaenkler/ddns-updater/pkg/server/totps"
)

func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Errorf("[api-UpdateConfig-1] could not parse form err: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	port, err := strconv.ParseUint(r.FormValue("port"), 10, 16)
	if err != nil {
		log.Errorf("[api-UpdateConfig-2] port is not valid: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	interval, err := strconv.ParseUint(r.FormValue("interval"), 10, 64)
	if err != nil {
		log.Errorf("[api-UpdateConfig-3] interval is not valid: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if interval < 10 {
		log.Errorf("[api-UpdateConfig-4] interval is too small: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resolver := strings.TrimSpace(r.FormValue("resolver"))
	if resolver != "" {
		_, err = url.ParseRequestURI(resolver)
		if err != nil {
			log.Errorf("[api-UpdateConfig-5] resolver is not valid: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	cfg := &config.Config{
		Port:     port,
		Interval: interval,
		Resolver: resolver,
	}
	if totps.Verify(r.FormValue("otp")) {
		cfg.UseTOTP = !config.Get().UseTOTP
	}
	config.Update(cfg)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusTemporaryRedirect)
}
