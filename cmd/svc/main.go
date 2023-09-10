package main

import (
	"os"
	"os/signal"

	"github.com/kardianos/service"
	"github.com/plaenkler/ddns-updater/pkg/database"
	"github.com/plaenkler/ddns-updater/pkg/ddns"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"github.com/plaenkler/ddns-updater/pkg/server"
	"github.com/plaenkler/ddns-updater/pkg/server/session"
)

type program struct{}

func (p *program) Start(_ service.Service) error {
	database.Start()
	log.Infof("[main-Start-1] started database connection")
	go ddns.Start()
	log.Infof("[main-Start-2] started ddns service")
	go session.Start()
	log.Infof("[main-Start-3] started session service")
	go server.Start()
	log.Infof("[main-Start-4] started webserver")
	return nil
}

func (p *program) Stop(_ service.Service) error {
	server.Stop()
	log.Infof("[main-Stop-1] stopped webserver")
	session.Stop()
	log.Infof("[main-Stop-2] stopped session service")
	ddns.Stop()
	log.Infof("[main-Stop-3] stopped ddns service")
	database.Stop()
	log.Infof("[main-Stop-4] stopped database connection")
	return nil
}

func main() {
	if service.AvailableSystems()[len(service.AvailableSystems())-1].Interactive() {
		p := &program{}
		p.Start(nil)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		p.Stop(nil)
		return
	}
	svcConfig := &service.Config{
		Name:        "DDNS-Updater",
		DisplayName: "Dynamic DNS Updater",
		Description: "Service to update dynamic DNS entries.",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalf("[main-main-1] failed to create service: %v", err)
	}
	err = s.Run()
	if err != nil {
		log.Errorf("[main-main-2] failed to run service: %v", err)
	}
}
