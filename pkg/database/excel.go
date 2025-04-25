package database

import (
	"fmt"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/xuri/excelize/v2"
)

func (q *Query) DownloadExcel(req *models.ExcelDownloadRequest) (*excelize.File, error) {
	query := `
	SELECT 
		f.student_unit_id AS usn,
		COALESCE(s.student_name, 'N/A') AS name,
		CASE
			WHEN a.login IS NOT NULL THEN 'P'
			ELSE 'A'
		END AS attendance_status
	FROM fingerprintdata f
	LEFT JOIN attendance a 
		ON f.student_id = a.student_id 
		AND a.unit_id = $1 
		AND a.date BETWEEN $2 AND $3
	LEFT JOIN student_details s 
		ON s.student_unit_id = f.student_unit_id
	WHERE f.unit_id = $1
	ORDER BY name;
	`

	rows, err := q.db.Query(query, req.UnitId, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	headers := []string{"Name", "USN", "Attendance Status"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	rowNum := 2
	for rows.Next() {
		var info models.ExcelStudentInfo
		var status string
		if err := rows.Scan(&info.Usn, &info.Name, &status); err != nil {
			return nil, err
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), info.Name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), info.Usn)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), status)
		rowNum++
	}

	return f, nil
}
