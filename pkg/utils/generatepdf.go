package utils

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/jung-kurt/gofpdf"
)

func formatTo12Hour(timeStr string) string {
	if timeStr == "pending" || timeStr == "Pending"{
		return timeStr
	}
	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return timeStr
	}
	return t.Format("03:04 PM")
}

func GeneratePDF(students []models.Student) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add header image
	headerImagePath := "pkg/utils/header.png"
	pdf.ImageOptions(headerImagePath, 10, 10, 190, 25, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	pdf.Ln(25)

	// Title
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Attendance Sheet")
	pdf.Ln(10)

	// Table Headers
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(20, 10, "S.No", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 10, "Student Name", "1", 0, "C", true, 0, "")
	pdf.CellFormat(50, 10, "USN", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 10, "Presence", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Table Data
	pdf.SetFont("Arial", "", 12)
	for i, student := range students {
		pdf.CellFormat(20, 10, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 10, student.StudentName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 10, student.StudentUSN, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, formatTo12Hour(student.LoginTime), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, formatTo12Hour(student.LogoutTime), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("Error while generating PDF: %v", err)
		return nil, err
	}

	log.Println("PDF generated successfully")
	return &buf, nil
}
