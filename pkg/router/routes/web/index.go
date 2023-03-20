package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/plaenkler/ddns/pkg/config"
	"github.com/plaenkler/ddns/pkg/database"
	"github.com/plaenkler/ddns/pkg/database/model"
	"github.com/plaenkler/ddns/pkg/ddns"
)

var (
	//go:embed static
	static embed.FS
)

type structIndex struct {
	Jobs      []model.SyncJob
	IPAddress string
	Config    *config.Config
	Providers []string
}

func ProvideIndex(writer http.ResponseWriter, request *http.Request) {
	// Default to 404
	if request.URL.Path != "/" {
		fmt.Fprintf(writer, "[provide-index-1] 404 - Not found")
		return
	}
	template, err := template.New("index").ParseFS(static,
		"static/html/pages/index.html",
		"static/html/partials/include.html",
		"static/html/partials/navbar.html",
		"static/html/partials/modals.html",
	)
	if err != nil {
		fmt.Fprintf(writer, "[provide-index-2] could not provide template - error: %s", err)
		return
	}
	structIndex := structIndex{
		Config:    config.GetConfig(),
		Providers: ddns.GetProviders(),
	}
	structIndex.IPAddress, err = ddns.GetPublicIP()
	if err != nil {
		fmt.Fprintf(writer, "[provide-index-3] could not get public IP address - error: %s", err)
		return
	}
	err = database.GetManager().DB.Find(&structIndex.Jobs).Error
	if err != nil {
		fmt.Fprintf(writer, "[provide-index-4] could not find jobs - error: %s", err)
		return
	}
	writer.Header().Add("Content-Type", "text/html")
	err = template.Execute(writer, structIndex)
	if err != nil {
		fmt.Fprintf(writer, "[provide-index-5] could not execute parsed template - error: %v", err)
	}
}
