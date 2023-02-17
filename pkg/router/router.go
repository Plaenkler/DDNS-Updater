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
	"github.com/plaenkler/ddns/pkg/router/routes"
)

var (
	//go:embed routes/static
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

		manager.Router.HandleFunc("/",
			routes.ProvideIndex)

		err := manager.provideFiles()
		if err != nil {
			log.Fatalf("[router-start-1] could not provide files - error: %s", err)
		}
		config := config.GetConfig()
		server := &http.Server{
			Addr:              fmt.Sprintf(":%v", config.Port),
			ReadTimeout:       3 * time.Second,
			ReadHeaderTimeout: 3 * time.Second,
			WriteTimeout:      3 * time.Second,
			IdleTimeout:       120 * time.Second,
			Handler:           manager.Router,
		}
		err = server.ListenAndServe()
		if err != nil {
			log.Fatalf("[router-start-2] failed starting router - error: %s", err.Error())
		}
	})
}

func (manager *Manager) provideFiles() error {
	fs, err := fs.Sub(static, "routes/static")
	if err != nil {
		return err
	}
	manager.Router.Handle("/js/", http.FileServer(http.FS(fs)))
	manager.Router.Handle("/css/", http.FileServer(http.FS(fs)))
	manager.Router.Handle("/img/", http.FileServer(http.FS(fs)))
	return nil
}
