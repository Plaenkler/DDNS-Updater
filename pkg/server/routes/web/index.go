package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"html"
	"github.com/plaenkler/ddns-updater/pkg/cipher"
	"github.com/plaenkler/ddns-updater/pkg/config"
	"github.com/plaenkler/ddns-updater/pkg/database"
	"github.com/plaenkler/ddns-updater/pkg/database/model"
	"github.com/plaenkler/ddns-updater/pkg/ddns"
	"github.com/plaenkler/ddns-updater/pkg/server/totps"
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
	addr, err := ddns.GetPublicIP()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not get public IP address: %s", html.EscapeString(err.Error()))
		return
	}
	img, err := totps.GetKeyAsQR()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not generate TOTP QR code: %s", html.EscapeString(err.Error()))
		return
	}
	data := indexPageData{
		Config:    config.Get(),
		Providers: ddns.GetProviders(),
		IPAddress: addr,
		TOTPQR:    template.URL(img),
	}
	db := database.GetDatabase()
	if db == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not get database connection")
		return
	}
	err = db.Find(&data.Jobs).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not find jobs: %s", html.EscapeString(err.Error()))
		return
	}
	err = sanitizeParams(data.Jobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "formatting params failed: %s", html.EscapeString(err.Error()))
		return
	}
	tpl, err := template.New("index").Funcs(template.FuncMap{
		"formatParams": formatParams,
	}).ParseFS(static,
		"static/html/pages/index.html",
		"static/html/partials/include.html",
		"static/html/partials/navbar.html",
		"static/html/partials/modals.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not provide template: %s", html.EscapeString(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "text/html")
	err = tpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not execute parsed template: %v", html.EscapeString(err.Error()))
	}
}

func sanitizeParams(jobs []model.SyncJob) error {
	for j := range jobs {
		decParams, err := cipher.Decrypt(jobs[j].Params)
		if err != nil {
			return err
		}
		params := make(map[string]string)
		err = json.Unmarshal(decParams, &params)
		if err != nil {
			return err
		}
		for k := range params {
			kl := strings.ToLower(k)
			if strings.Contains(kl, "password") {
				params[k] = "***"
			}
			if strings.Contains(kl, "token") {
				params[k] = "***"
			}
		}
		encParams, err := json.Marshal(params)
		if err != nil {
			return err
		}
		jobs[j].Params = string(encParams)
	}
	return nil
}

func formatParams(paramsData string) (string, error) {
	params := make(map[string]string)
	err := json.Unmarshal([]byte(paramsData), &params)
	if err != nil {
		return "", err
	}
	var formatted string
	for key, value := range params {
		formatted += fmt.Sprintf("%s: %v, ", key, value)
	}
	formatted = formatted[:len(formatted)-2]
	return formatted, nil
}
