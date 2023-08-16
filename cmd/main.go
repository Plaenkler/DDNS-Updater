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
	log.Infof("[main-main-1] started database connection")
	go ddns.Start()
	log.Infof("[main-main-2] started ddns service")
	go session.Start()
	log.Infof("[main-main-3] started session service")
	go server.Start()
	log.Infof("[main-main-4] started webserver")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	server.Stop()
	log.Infof("[main-main-5] stopped webserver")
	session.Stop()
	log.Infof("[main-main-6] stopped session service")
	ddns.Stop()
	log.Infof("[main-main-7] stopped ddns service")
	database.Stop()
	log.Infof("[main-main-8] stopped database connection")
}
