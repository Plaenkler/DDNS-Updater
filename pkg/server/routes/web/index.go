package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/plaenkler/ddns-updater/pkg/config"
	"github.com/plaenkler/ddns-updater/pkg/database"
	"github.com/plaenkler/ddns-updater/pkg/database/model"
	"github.com/plaenkler/ddns-updater/pkg/ddns"
)

var (
	//go:embed static
	static embed.FS
)

type templateData struct {
	Jobs      []model.SyncJob
	IPAddress string
	Config    *config.Config
	Providers []string
}

func ProvideIndex(w http.ResponseWriter, r *http.Request) {
	// Default to 404
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "[web-ProvideIndex-1] 404 - Not found")
		return
	}
	template, err := template.New("index").ParseFS(static,
		"static/html/pages/index.html",
		"static/html/partials/include.html",
		"static/html/partials/navbar.html",
		"static/html/partials/modals.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideIndex-2] could not provide template - error: %s", err)
		return
	}
	data := templateData{
		Config:    config.GetConfig(),
		Providers: ddns.GetProviders(),
	}
	data.IPAddress, err = ddns.GetPublicIP()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideIndex-3] could not get public IP address - error: %s", err)
		return
	}
	db := database.GetDatabase()
	if db == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideIndex-4] could not get database connection")
		return
	}
	err = db.Find(&data.Jobs).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideIndex-5] could not find jobs - error: %s", err)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	err = template.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideIndex-6] could not execute parsed template - error: %v", err)
	}
}
