package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/repository"
)

type ExcelController struct {
	Repo *repository.ExcelRepository
}

func NewExcelController(repo *repository.ExcelRepository) *ExcelController {
	return &ExcelController{Repo: repo}
}

// Helper function to parse different date formats
func parseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		time.RFC3339, // "2025-03-22T09:00:00Z"
		"02-01-2006", // "22-03-2025"
		"2006-01-02", // "2025-03-22"
		"02/01/2006", // "22/03/2025"
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}

func (c *ExcelController) DownloadExcel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ExcelDownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if all required fields are present
	if req.UnitId == "" || req.UserId == "" || req.Slot == "" || req.StartDate == "" || req.EndDate == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Parse the start and end dates using the helper function
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		http.Error(w, "Invalid start_date format: "+err.Error(), http.StatusBadRequest)
		return
	}
	endDate, err := parseDate(req.EndDate)
	if err != nil {
		http.Error(w, "Invalid end_date format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Format dates for DB query
	req.StartDate = startDate.Format("2006-01-02")
	req.EndDate = endDate.Format("2006-01-02")

	// Generate the Excel file
	file, err := c.Repo.DownloadExcel(&req)
	if err != nil {
		http.Error(w, "Failed to generate Excel file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the file locally (for debugging or later use)
	saveDir := "/Users/bunny/Desktop/Vithsutra"
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create save directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileName := "student_" + time.Now().Format("2006_01_02_15_04_05") + ".xlsx"
	fullPath := fmt.Sprintf("%s/%s", saveDir, fileName)

	if err := file.SaveAs(fullPath); err != nil {
		http.Error(w, "Failed to save Excel file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Excel downloaded and saved successfully",
		"filename": fileName,
		"path":     fullPath,
	})
}
