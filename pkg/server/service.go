package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/plaenkler/ddns/pkg/config"
	"github.com/plaenkler/ddns/pkg/server/routes/api"
	"github.com/plaenkler/ddns/pkg/server/routes/web"
)

var (
	lock   sync.Mutex
	static embed.FS
	router *http.ServeMux
)

func StartService() {
	lock.Lock()
	defer lock.Unlock()
	err := initializeRouter()
	if err != nil {
		log.Fatalf("[server-StartService-1] could not initialize router: %v", err)
	}
	err = initializeServer()
	if err != nil {
		log.Fatalf("[server-StartService-2] could not initialize server: %v", err)
	}
}

func initializeRouter() error {
	r := NewRouter()
	registerMiddlewares(r)
	registerAPIRoutes(r)
	router = r.ServeMux
	err := registerStaticFiles(router)
	if err != nil {
		return err
	}
	return nil
}

func registerMiddlewares(r *Router) {
	r.Use(forwardToProxy)
	r.Use(limitRequests)
}

func registerAPIRoutes(r *Router) {
	r.HandleFunc("/", web.ProvideIndex)
	r.HandleFunc("/api/inputs", api.GetInputs)
	r.HandleFunc("/api/job/create", api.CreateJob)
	r.HandleFunc("/api/job/update", api.UpdateJob)
	r.HandleFunc("/api/job/delete", api.DeleteJob)
	r.HandleFunc("/api/config/update", api.UpdateConfig)
}

func registerStaticFiles(router *http.ServeMux) error {
	fs, err := fs.Sub(static, "routes/web/static")
	if err != nil {
		return err
	}
	router.Handle("/js/", http.FileServer(http.FS(fs)))
	router.Handle("/css/", http.FileServer(http.FS(fs)))
	router.Handle("/img/", http.FileServer(http.FS(fs)))
	return nil
}

func initializeServer() error {
	config := config.GetConfig()
	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", config.Port),
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
		Handler:           router,
	}
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
