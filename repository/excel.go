package repository

import (
	"database/sql"
	"net/http"
)

type ExcelRepo struct {
	db*sql.DB
}

func NewExcelRepo(db *sql.DB) *ExcelRepo {
	return &ExcelRepo{
		db,
	}
}

func (er *ExcelRepo) GenerateExcelReport(r *http.Request) () {

	// Parse rows into attendance records
	records := make(map[string]*AttendanceRecord) // Key: USN
	for rows.Next() {
		var name, usn, day, status string
		if err := rows.Scan(&name, &usn, &day, &status); err != nil {
			log.Fatal(err)
		}

		// Initialize the record if not already present
		if _, exists := records[usn]; !exists {
			records[usn] = &AttendanceRecord{
				Name:       name,
				USN:        usn,
				Attendance: make(map[string]string),
			}
		}

		// Add the attendance for the specific day
		records[usn].Attendance[day] = status
	}

	// Generate Excel Sheet
	f := excelize.NewFile()
	sheetName := "Attendance"
	f.SetSheetName(f.GetSheetName(0), sheetName)

	// Header Row
	headers := []string{"Name", "USN"}
	for d := 1; d <= 31; d++ {
		headers = append(headers, fmt.Sprintf("%02d", d)) // Add day columns (01, 02, ..., 31)
	}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string('A'+i))
		f.SetCellValue(sheetName, cell, header)
	}

	// Write Records
	rowIndex := 2
	for _, record := range records {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), record.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), record.USN)

		// Fill attendance for each day
		for d := 1; d <= 31; d++ {
			day := fmt.Sprintf("%02d", d)
			status := record.Attendance[day]
			if status == "" {
				status = "A" // Mark absent if no record exists for the day
			}
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", string('A'+2+d-1), rowIndex), status)
		}
		rowIndex++
	}

	// Save the file
	if err := f.SaveAs("FormalAttendanceSheet.xlsx"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Attendance sheet generated: FormalAttendanceSheet.xlsx")
}
