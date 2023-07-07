package main

import (
	"log"

	"github.com/plaenkler/ddns/pkg/database"
	"github.com/plaenkler/ddns/pkg/ddns"
	"github.com/plaenkler/ddns/pkg/server"
)

func init() {
	database.GetManager().Start()
	log.Printf("[main-init-1] initialized database")
}

func main() {
	go ddns.Run()
	log.Printf("[main-main-1] started ddns service")
	server.StartService()
}
