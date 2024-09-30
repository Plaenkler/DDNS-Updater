package main

import (
	"os"
	"os/signal"

	"github.com/plaenkler/ddns-updater/pkg/database"
	"github.com/plaenkler/ddns-updater/pkg/ddns"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"github.com/plaenkler/ddns-updater/pkg/server"
	"github.com/plaenkler/ddns-updater/pkg/server/session"
)

func main() {
	database.Start()
	log.Infof("started database connection")
	go ddns.Start()
	log.Infof("started ddns service")
	go session.Start()
	log.Infof("started session service")
	go server.Start()
	log.Infof("started webserver")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	server.Stop()
	log.Infof("stopped webserver")
	session.Stop()
	log.Infof("stopped session service")
	ddns.Stop()
	log.Infof("stopped ddns service")
	database.Stop()
	log.Infof("stopped database connection")
}
