package database

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/plaenkler/ddns-updater/pkg/database/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	pathToDB = "./data/ddns.db"
)

var (
	managerOnce sync.Once
	startOnce   sync.Once
	instance    *Manager
)

type Manager struct {
	DB *gorm.DB
}

func GetManager() *Manager {
	managerOnce.Do(func() {
		instance = &Manager{}
	})
	return instance
}

func (manager *Manager) Start() {
	startOnce.Do(func() {
		db, err := manager.connect()
		if err != nil {
			log.Fatalf("[database-Start-1] connection failed - error: %s", err.Error())
		}
		manager.DB = db
		err = manager.DB.AutoMigrate(
			&model.SyncJob{},
			&model.IPAddress{},
		)
		if err != nil {
			log.Fatalf("[database-Start-2] migration failed - error: %s", err.Error())
		}
	})
}

func (m *Manager) connect() (*gorm.DB, error) {
	err := os.MkdirAll(filepath.Dir(pathToDB), 0755)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(pathToDB))
	if err != nil {
		return nil, err
	}
	return db, nil
}
