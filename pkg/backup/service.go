package backup

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/plaenkler/ddns-updater/pkg/cipher"
	"github.com/plaenkler/ddns-updater/pkg/config"
	"github.com/plaenkler/ddns-updater/pkg/database"
	"github.com/plaenkler/ddns-updater/pkg/database/model"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
)

// BackupData represents the complete backup structure
type BackupData struct {
	Version   string             `json:"version"`
	Timestamp time.Time          `json:"timestamp"`
	Config    *config.Config     `json:"config"`
	Jobs      []BackupJob        `json:"jobs"`
}

// BackupJob represents a job in the backup with decrypted parameters
type BackupJob struct {
	Provider string `json:"provider"`
	Params   string `json:"params"` // decrypted
}

// Export creates a complete backup of the system configuration and jobs
func Export() (*BackupData, error) {
	// Get current configuration
	cfg := config.Get()
	if cfg == nil {
		log.Errorf("could not get configuration")
		return nil, fmt.Errorf("could not get configuration")
	}

	// Get database connection
	db := database.GetDatabase()
	if db == nil {
		log.Errorf("could not get database connection")
		return nil, fmt.Errorf("could not get database connection")
	}

	// Retrieve all jobs from database
	var jobs []model.SyncJob
	if err := db.Find(&jobs).Error; err != nil {
		log.Errorf("could not retrieve jobs from database: %s", err.Error())
		return nil, fmt.Errorf("could not retrieve jobs from database: %s", err.Error())
	}

	// Decrypt job parameters for backup
	var backupJobs []BackupJob
	for _, job := range jobs {
		decryptedParams, err := cipher.Decrypt(job.Params)
		if err != nil {
			log.Errorf("could not decrypt params for job %d: %s", job.ID, err.Error())
			continue
		}

		backupJobs = append(backupJobs, BackupJob{
			Provider: job.Provider,
			Params:   string(decryptedParams),
		})
	}

	backup := &BackupData{
		Version:   "1.0",
		Timestamp: time.Now(),
		Config:    cfg,
		Jobs:      backupJobs,
	}

	return backup, nil
}

// Import restores the system configuration and jobs from backup data
func Import(backupData *BackupData) error {
	if backupData.Version != "1.0" {
		log.Errorf("unsupported backup version: %s", backupData.Version)
		return fmt.Errorf("unsupported backup version: %s", backupData.Version)
	}

	// Get database connection
	db := database.GetDatabase()
	if db == nil {
		log.Errorf("could not get database connection")
		return fmt.Errorf("could not get database connection")
	}

	// Begin transaction
	tx := db.Begin()
	if tx.Error != nil {
		log.Errorf("could not begin transaction: %s", tx.Error.Error())
		return fmt.Errorf("could not begin transaction: %s", tx.Error.Error())
	}

	// Clear existing jobs
	if err := tx.Unscoped().Delete(&model.SyncJob{}, "1 = 1").Error; err != nil {
		tx.Rollback()
		log.Errorf("could not clear existing jobs: %s", err.Error())
		return fmt.Errorf("could not clear existing jobs: %s", err.Error())
	}

	// Restore jobs
	for _, backupJob := range backupData.Jobs {
		// Encrypt parameters for storage
		encryptedParams, err := cipher.Encrypt(backupJob.Params)
		if err != nil {
			log.Errorf("could not encrypt params for job with provider %s: %s", backupJob.Provider, err.Error())
			continue
		}

		job := model.SyncJob{
			Provider: backupJob.Provider,
			Params:   encryptedParams,
		}

		if err := tx.Create(&job).Error; err != nil {
			log.Errorf("could not create job with provider %s: %s", backupJob.Provider, err.Error())
			continue
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Errorf("could not commit transaction: %s", err.Error())
		return fmt.Errorf("could not commit transaction: %s", err.Error())
	}

	// Update configuration
	if err := config.Update(backupData.Config); err != nil {
		log.Errorf("could not update configuration: %s", err.Error())
		return fmt.Errorf("could not update configuration: %s", err.Error())
	}

	log.Infof("successfully imported backup with %d jobs", len(backupData.Jobs))
	return nil
}

// Marshal converts backup data to JSON
func (b *BackupData) Marshal() ([]byte, error) {
	return json.MarshalIndent(b, "", "  ")
}

// UnmarshalBackupData parses JSON backup data
func UnmarshalBackupData(data []byte) (*BackupData, error) {
	var backup BackupData
	if err := json.Unmarshal(data, &backup); err != nil {
		return nil, err
	}
	return &backup, nil
}