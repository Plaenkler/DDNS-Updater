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
	"github.com/plaenkler/ddns-updater/pkg/server/totp"
)

var (
	//go:embed static
	static embed.FS
)

type indexPageData struct {
	Jobs      []model.SyncJob
	IPAddress string
	Config    *config.Config
	Providers []string
	TOTPQR    template.URL
}

func ProvideIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tpl, err := template.New("index").ParseFS(static,
		"static/html/pages/index.html",
		"static/html/partials/include.html",
		"static/html/partials/navbar.html",
		"static/html/partials/modals.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideIndex-1] could not provide template: %s", err)
		return
	}
	data := indexPageData{
		Config:    config.GetConfig(),
		Providers: ddns.GetProviders(),
	}
	data.IPAddress, err = ddns.GetPublicIP()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideIndex-2] could not get public IP address: %s", err)
		return
	}
	db := database.GetDatabase()
	if db == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideIndex-3] could not get database connection")
		return
	}
	err = db.Find(&data.Jobs).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideIndex-4] could not find jobs: %s", err)
		return
	}
	img, err := totp.GetKeyAsQR()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideIndex-5] could not generate TOTP QR code: %s", err)
		return
	}
	data.TOTPQR = template.URL(img)
	w.Header().Add("Content-Type", "text/html")
	err = tpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "[web-ProvideIndex-6] could not execute parsed template: %v", err)
	}
}
