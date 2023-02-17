package database

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/plaenkler/ddns/pkg/model"
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
		db, err := instance.connect()
		if err != nil {
			log.Fatalf("[start-1] connection failed - error: %s", err.Error())
		}
		manager.DB = db
		err = manager.DB.AutoMigrate(
			&model.User{},
			&model.Updater{},
		)
		if err != nil {
			log.Fatalf("[start-2] migration failed - error: %s", err.Error())
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
