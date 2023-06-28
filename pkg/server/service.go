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
		instance = &Manager{
			Router: http.NewServeMux(),
		}
	})
	return instance
}

func (m *Manager) Start() {
	startOnce.Do(func() {
		m.initializeRouter()
		err := m.provideFiles()
		if err != nil {
			log.Fatalf("[server-Start-1] could not provide files - error: %s", err)
		}
		m.initializeServer()
	})
}

func (m *Manager) initializeRouter() {
	r := NewRouter()
	r.Use(limitRequests)
	r.HandleFunc("/", web.ProvideIndex)
	r.HandleFunc("/api/inputs", api.GetInputs)
	r.HandleFunc("/api/job/create", api.CreateJob)
	r.HandleFunc("/api/job/update", api.UpdateJob)
	r.HandleFunc("/api/job/delete", api.DeleteJob)
	r.HandleFunc("/api/config/update", api.UpdateConfig)
	m.Router = r.ServeMux
}

func (m *Manager) initializeServer() {
	config := config.GetConfig()
	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", config.Port),
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read X-Forwarded-For-Header for proxy requests
			remoteIP := r.Header.Get("X-Forwarded-For")
			remotePort := r.Header.Get("X-Forwarded-Port")
			if remotePort == "" {
				remotePort = "80"
			}
			r.URL.Scheme = "http"
			r.URL.Host = fmt.Sprintf("%s:%s", remoteIP, remotePort)
			m.Router.ServeHTTP(w, r)
		}),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("[server-initializeServer-1] failed starting router - error: %s", err.Error())
	}
}

func (m *Manager) provideFiles() error {
	fs, err := fs.Sub(static, "routes/web/static")
	if err != nil {
		return err
	}
	m.Router.Handle("/js/", http.FileServer(http.FS(fs)))
	m.Router.Handle("/css/", http.FileServer(http.FS(fs)))
	m.Router.Handle("/img/", http.FileServer(http.FS(fs)))
	return nil
}

func limitRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := IsOverLimit(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
