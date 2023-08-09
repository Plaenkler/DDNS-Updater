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
	database.StartService()
	log.Infof("[main-main-1] started database connection")
	go ddns.StartService()
	log.Infof("[main-main-2] started ddns service")
	go session.StartService()
	log.Infof("[main-main-3] started session service")
	go server.StartService()
	log.Infof("[main-main-4] started webserver")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	server.StopService()
	log.Infof("[main-main-5] stopped webserver")
	session.StopService()
	log.Infof("[main-main-6] stopped session service")
	ddns.StopService()
	log.Infof("[main-main-7] stopped ddns service")
	database.StopService()
	log.Infof("[main-main-8] stopped database connection")
}
