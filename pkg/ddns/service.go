package ddns

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/plaenkler/ddns-updater/pkg/config"
	"github.com/plaenkler/ddns-updater/pkg/database"
	"github.com/plaenkler/ddns-updater/pkg/database/model"
)

var run sync.Once

func Run() {
	run.Do(func() {
		ticker := time.NewTicker(time.Second * time.Duration(config.GetConfig().Interval))
		defer ticker.Stop()
		for range ticker.C {
			address, err := GetPublicIP()
			if err != nil {
				log.Printf("[ddns-Run-1] failed to get public IP address - error: %v", err)
				continue
			}
			newAddress := model.IPAddress{
				Address: address,
			}
			err = database.GetManager().DB.FirstOrCreate(&newAddress, newAddress).Error
			if err != nil {
				log.Printf("[ddns-Run-2] failed to save new IP address - error: %v", err)
				continue
			}
			jobs := []model.SyncJob{}
			err = database.GetManager().DB.Where("NOT ip_address_id = ? or ip_address_id IS NULL", newAddress.ID).Find(&jobs).Error
			if err != nil {
				log.Printf("[ddns-Run-3] failed to get DDNS update jobs - error: %v", err)
				continue
			}
			if len(jobs) == 0 {
				log.Println("[ddns-Run-4] no dynamic DNS record needs to be updated")
				continue
			}
			for _, job := range jobs {
				updater, ok := updaters[job.Provider]
				if !ok {
					log.Printf("[ddns-Run-5] no updater found for job %v", job.ID)
					continue
				}
				request := updater.Request
				err := json.Unmarshal([]byte(job.Params), request)
				if err != nil {
					log.Printf("[ddns-Run-6] failed to unmarshal job params for job %v - error: %s", job.ID, err)
					continue
				}
				err = updater.Updater(request, address)
				if err != nil {
					log.Printf("[ddns-Run-7] failed to update DDNS entry for job %v - error: %s", job.ID, err)
					continue
				}
				err = database.GetManager().DB.Model(&job).Update("ip_address_id", newAddress.ID).Error
				if err != nil {
					log.Printf("[ddns-Run-8] failed to update IP address for job %v - error: %s", job.ID, err)
				}
				log.Printf("[ddns-Run-9] updated DDNS entry for ID: %v Provider: %s Params: %+v", job.ID, job.Provider, job.Params)
			}
		}
	})
}
