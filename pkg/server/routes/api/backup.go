package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/plaenkler/ddns-updater/pkg/backup"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
)

// ExportBackup handles the backup export API endpoint
func ExportBackup(w http.ResponseWriter, r *http.Request) {
	// Export backup data
	backupData, err := backup.Export()
	if err != nil {
		log.Errorf("could not export backup: %s", err.Error())
		http.Error(w, "Could not export backup", http.StatusInternalServerError)
		return
	}

	// Convert to JSON
	jsonData, err := backupData.Marshal()
	if err != nil {
		log.Errorf("could not marshal backup data: %s", err.Error())
		http.Error(w, "Could not marshal backup data", http.StatusInternalServerError)
		return
	}

	// Set headers for file download
	filename := fmt.Sprintf("ddns-backup-%s.json", backupData.Timestamp.Format("2006-01-02-15-04-05"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(jsonData)))

	// Write JSON data
	_, err = w.Write(jsonData)
	if err != nil {
		log.Errorf("could not write backup data: %s", err.Error())
		return
	}

	log.Infof("exported backup with %d jobs", len(backupData.Jobs))
}

// ImportBackup handles the backup import API endpoint
func ImportBackup(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		log.Errorf("could not parse multipart form: %s", err.Error())
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, _, err := r.FormFile("backup")
	if err != nil {
		log.Errorf("could not get uploaded file: %s", err.Error())
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Errorf("could not read file content: %s", err.Error())
		http.Error(w, "Could not read file content", http.StatusInternalServerError)
		return
	}

	// Parse backup data
	backupData, err := backup.UnmarshalBackupData(fileContent)
	if err != nil {
		log.Errorf("could not parse backup data: %s", err.Error())
		http.Error(w, "Invalid backup file format", http.StatusBadRequest)
		return
	}

	// Import backup data
	err = backup.Import(backupData)
	if err != nil {
		log.Errorf("could not import backup: %s", err.Error())
		http.Error(w, "Could not import backup", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "Backup imported successfully"}`))

	log.Infof("imported backup with %d jobs", len(backupData.Jobs))
}