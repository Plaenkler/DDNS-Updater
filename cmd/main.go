package main

import (
	"github.com/plaenkler/ddns/pkg/database"
	"github.com/plaenkler/ddns/pkg/ddns"
	"github.com/plaenkler/ddns/pkg/router"
)

func init() {
	database.GetManager().Start()
}

func main() {
	go ddns.Start()
	router.GetManager().Start()
}
