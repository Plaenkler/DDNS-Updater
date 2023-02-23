package router

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/plaenkler/ddns/pkg/config"
	"github.com/plaenkler/ddns/pkg/router/routes/api"
	"github.com/plaenkler/ddns/pkg/router/routes/web"
)

var (
	//go:embed routes/web/static
	static      embed.FS
	managerOnce sync.Once
	startOnce   sync.Once
	instance    *Manager
)

type Manager struct {
	Router *http.ServeMux
}

func GetManager() *Manager {
	managerOnce.Do(func() {
		instance = &Manager{}
	})
	return instance
}

func (manager *Manager) Start() {
	startOnce.Do(func() {
		manager.Router = http.NewServeMux()
		manager.Router.HandleFunc("/", web.ProvideIndex)
		manager.Router.HandleFunc("/api/inputs", api.GetInputs)
		manager.Router.HandleFunc("/api/job/create", api.CreateJob)
		manager.Router.HandleFunc("/api/job/update", api.UpdateJob)
		manager.Router.HandleFunc("/api/job/delete", api.DeleteJob)
		manager.Router.HandleFunc("/api/config/update", api.UpdateConfig)
		err := manager.provideFiles()
		if err != nil {
			log.Fatalf("[router-start-1] could not provide files - error: %s", err)
		}
		config := config.GetConfig()
		server := &http.Server{
			Addr:              fmt.Sprintf(":%v", config.Port),
			ReadTimeout:       15 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       120 * time.Second,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Read X-Forwarded-For-Header for proxying
				remoteIP := r.Header.Get("X-Forwarded-For")
				remotePort := r.Header.Get("X-Forwarded-Port")
				if remotePort == "" {
					remotePort = "80"
				}
				r.URL.Scheme = "http"
				r.URL.Host = fmt.Sprintf("%s:%s", remoteIP, remotePort)
				manager.Router.ServeHTTP(w, r)
			}),
		}
		err = server.ListenAndServe()
		if err != nil {
			log.Fatalf("[router-start-2] failed starting router - error: %s", err.Error())
		}
	})
}

func (manager *Manager) provideFiles() error {
	fs, err := fs.Sub(static, "routes/web/static")
	if err != nil {
		return err
	}
	manager.Router.Handle("/js/", http.FileServer(http.FS(fs)))
	manager.Router.Handle("/css/", http.FileServer(http.FS(fs)))
	manager.Router.Handle("/img/", http.FileServer(http.FS(fs)))
	return nil
}
