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
	log.Infof("started database connection")
	go ddns.Start()
	log.Infof("started ddns service")
	go session.Start()
	log.Infof("started session service")
	go server.Start()
	log.Infof("started webserver")
	return nil
}

func (p *program) Stop(_ service.Service) error {
	server.Stop()
	log.Infof("stopped webserver")
	session.Stop()
	log.Infof("stopped session service")
	ddns.Stop()
	log.Infof("stopped ddns service")
	database.Stop()
	log.Infof("stopped database connection")
	return nil
}

func main() {
	if service.AvailableSystems()[len(service.AvailableSystems())-1].Interactive() {
		p := &program{}
		err := p.Start(nil)
		if err != nil {
			log.Fatalf("failed to start service: %v", err)
		}
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		err = p.Stop(nil)
		if err != nil {
			log.Fatalf("failed to stop service: %v", err)
		}
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
		log.Fatalf("failed to create service: %v", err)
	}
	err = s.Run()
	if err != nil {
		log.Errorf("failed to run service: %v", err)
	}
}
