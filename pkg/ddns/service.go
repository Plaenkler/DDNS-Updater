package ddns

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/plaenkler/ddns/pkg/database"
	"github.com/plaenkler/ddns/pkg/model"
)

type ResolverFunc func(job model.Updater, ipAddr string) error

func Start() {
	// Register all DDNS providers
	resolvers := map[string]ResolverFunc{
		"Strato": updateStrato,
		"DDNSS":  updateDDNSS,
	}
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		var jobs []model.Updater
		err := database.GetManager().DB.Find(&jobs).Error
		if err != nil {
			log.Printf("[service-start-1] failed to get DDNS update jobs: %v\n", err)
			continue
		}
		ipAddr, err := getPublicIP()
		if err != nil {
			log.Printf("[service-start-2] failed to get public IP address: %v\n", err)
			continue
		}
		for _, job := range jobs {
			if job.LastIPAddr == ipAddr {
				log.Printf("[service-start-3] IP address for %q has not changed\n", job.Domain)
				continue
			}
			resolver, ok := resolvers[job.Provider]
			if !ok {
				log.Printf("[service-start-4] no resolver found for provider %q\n", job.Provider)
				continue
			}
			err = resolver(job, ipAddr)
			if err != nil {
				log.Printf("[service-start-5] failed to update DDNS entry for %q: %v\n", job.Domain, err)
			}
		}
	}
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ipBytes), nil
}
