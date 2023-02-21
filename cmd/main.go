package main

import (
	"log"

	"github.com/plaenkler/ddns/pkg/database"
	"github.com/plaenkler/ddns/pkg/ddns"
	"github.com/plaenkler/ddns/pkg/router"
)

func init() {
	database.GetManager().Start()
	log.Printf("[main-init-1] database started")
}

func main() {
	go ddns.Run()
	log.Printf("[main-main-1] ddns service started")
	router.GetManager().Start()
}
