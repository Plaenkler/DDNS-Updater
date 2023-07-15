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
	start()
	stop()
}

func start() {
	database.StartService()
	log.Infof("[main-start-1] started database connection")
	go ddns.StartService()
	log.Infof("[main-start-2] started ddns service")
	go server.StartService()
	log.Infof("[main-start-3] started webserver")
}

func stop() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	server.StopService()
	log.Infof("[main-stop-1] stopped webserver")
	ddns.StopService()
	log.Infof("[main-stop-2] stopped ddns service")
	database.StopService()
	log.Infof("[main-stop-3] stopped database connection")
}
