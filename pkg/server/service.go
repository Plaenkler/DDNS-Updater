package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/plaenkler/ddns-updater/pkg/config"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"github.com/plaenkler/ddns-updater/pkg/server/routes/api"
	"github.com/plaenkler/ddns-updater/pkg/server/routes/web"
)

var (
	//go:embed routes/web/static
	static embed.FS
	oc     sync.Once
	router *http.ServeMux
	server *http.Server
)

func StartService() {
	oc.Do(func() {
		initializeRouter()
		initializeServer()
	})
}

func initializeRouter() {
	r := NewRouter()
	registerMiddlewares(r)
	registerAPIRoutes(r)
	registerStaticFiles(r)
	router = r.ServeMux
}

func registerMiddlewares(r *Router) {
	r.Use(forwardToProxy)
	r.Use(limitRequests)
}

func registerAPIRoutes(r *Router) {
	r.HandleFunc("/", web.ProvideIndex)
	r.HandleFunc("/login", web.ProvideLogin)
	r.HandleFunc("/api/inputs", api.GetInputs)
	r.HandleFunc("/api/job/create", api.CreateJob)
	r.HandleFunc("/api/job/update", api.UpdateJob)
	r.HandleFunc("/api/job/delete", api.DeleteJob)
	r.HandleFunc("/api/config/update", api.UpdateConfig)
}

func registerStaticFiles(r *Router) {
	staticHandler := createStaticHandler()
	r.Handle("/js/", staticHandler)
	r.Handle("/css/", staticHandler)
}

func createStaticHandler() http.Handler {
	fs, err := fs.Sub(static, "routes/web/static")
	if err != nil {
		log.Fatalf("[server-createStaticHandler-1] could not create static handler: %v", err)
	}
	return controlCache(http.FileServer(http.FS(fs)))
}

func initializeServer() {
	server = &http.Server{
		Addr:              fmt.Sprintf(":%v", config.GetConfig().Port),
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
		Handler:           router,
	}
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("[server-initializeServer-1] could not initialize server: %v", err)
	}
}

func StopService() {
	if server == nil {
		return
	}
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Errorf("could not shutdown server: %v", err)
	}
}
