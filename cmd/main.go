package main

import (
	"github.com/plaenkler/ddns/pkg/database"
	"github.com/plaenkler/ddns/pkg/router"
)

func init() {
	database.GetManager().Start()
}

func main() {
	router.GetManager().Start()
}
