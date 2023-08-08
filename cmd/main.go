package main

import (
	"os"
	"os/signal"

	"github.com/plaenkler/ddns-updater/pkg/database"
	"github.com/plaenkler/ddns-updater/pkg/ddns"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"github.com/plaenkler/ddns-updater/pkg/server"
)

func main() {
	database.StartService()
	log.Infof("[main-main-1] started database connection")
	go ddns.StartService()
	log.Infof("[main-main-2] started ddns service")
	go server.StartService()
	log.Infof("[main-main-3] started webserver")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	server.StopService()
	log.Infof("[main-main-4] stopped webserver")
	ddns.StopService()
	log.Infof("[main-main-5] stopped ddns service")
	database.StopService()
	log.Infof("[main-main-6] stopped database connection")
}
