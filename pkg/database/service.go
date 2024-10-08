package database

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/plaenkler/ddns-updater/pkg/database/model"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	pathToDB = "./data/ddns.db"
)

var (
	db *gorm.DB
	oc sync.Once
)

func Start() {
	oc.Do(func() {
		err := createDBDir()
		if err != nil {
			log.Fatalf("failed to create database directory: %s", err.Error())
		}
		db, err = openDBConnection()
		if err != nil {
			log.Fatalf("failed to open database connection: %s", err.Error())
		}
		err = migrateDBSchema(db)
		if err != nil {
			log.Fatalf("failed to migrate database schema: %s", err.Error())
		}
	})
}

func createDBDir() error {
	err := os.MkdirAll(filepath.Dir(pathToDB), 0755)
	if err != nil {
		return err
	}
	return nil
}

func openDBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(pathToDB))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func migrateDBSchema(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.SyncJob{},
		&model.IPAddress{},
	)
	if err != nil {
		return err
	}
	return nil
}

func Stop() {
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("failed to get underlying DB connection: %s", err.Error())
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Errorf("failed to close DB connection: %s", err.Error())
	}
}

func GetDatabase() *gorm.DB {
	return db
}
