package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
)

type studentRepo struct {
	db *sql.DB
}

func NewStudentRepo(db *sql.DB) *studentRepo {
	return &studentRepo{
		db,
	}
}

func (repo *studentRepo) CreateNewStudent(r *http.Request) error {
	var createStudentRequest models.CreateStudentRequest

	if err := json.NewDecoder(r.Body).Decode(&createStudentRequest); err != nil {
		return errors.New("invalid json format")
	}

	validate := validator.New()
	if err := validate.Struct(createStudentRequest); err != nil {
		return errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	isStudentUnitIdExists, err := query.CheckStudentUnitIdExists(createStudentRequest.UnitId, createStudentRequest.StudentUnitId)

	if err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}

	if isStudentUnitIdExists {
		return errors.New("student unit id already exists")
	}

	var student models.Student

	unitId := strings.ToLower(createStudentRequest.UnitId)

	student.StudentId = uuid.NewString()
	student.StudentUnitId = createStudentRequest.StudentUnitId
	student.StudentName = createStudentRequest.StudentName
	student.StudentUsn = createStudentRequest.StudentUsn
	student.Department = createStudentRequest.Department

	if err := query.CreateNewStudent(&student, unitId, createStudentRequest.FingerprintData); err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}

	return nil
}

func (repo *studentRepo) UpdateStudentDetails(r *http.Request) error {
	var studentUpdateRequest models.UpdateStudentRequest

	if err := json.NewDecoder(r.Body).Decode(&studentUpdateRequest); err != nil {
		return errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(studentUpdateRequest); err != nil {
		return errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	unitId := strings.ToLower(studentUpdateRequest.UnitId)
	if err := query.UpdateStudent(unitId, studentUpdateRequest.StudentId, studentUpdateRequest.StudentName, studentUpdateRequest.StudentUsn, studentUpdateRequest.Department); err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}

	return nil
}

func (repo *studentRepo) DeleteStudent(r *http.Request) error {
	var deleteStudentRequest models.DeleteStudentRequest

	if err := json.NewDecoder(r.Body).Decode(&deleteStudentRequest); err != nil {
		return errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(deleteStudentRequest); err != nil {
		return errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	unitId := strings.ToLower(deleteStudentRequest.UnitId)

	if err := query.DeleteStudent(unitId, deleteStudentRequest.StudentId, deleteStudentRequest.StudentUnitId); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (repo *studentRepo) GetStudentDetails(r *http.Request) ([]*models.Student, error) {
	vars := mux.Vars(r)

	unitId := vars["unit_id"]

	query := database.NewQuery(repo.db)

	students, err := query.GetStudentDetails(unitId)

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}

	return students, nil
}

func (repo *studentRepo) GetStudentLogs(r *http.Request) ([]*models.StudentAttendanceLog, error) {
	studentId := mux.Vars(r)["student_id"]

	query := database.NewQuery(repo.db)

	logs, err := query.GetStudentLogs(studentId)

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}

	return logs, nil

}

func (repo *studentRepo) DownloadPdf(r *http.Request) (*gopdf.GoPdf, error) {
	var pdfDownloadRequest models.PdfDownloadRequest

	if err := json.NewDecoder(r.Body).Decode(&pdfDownloadRequest); err != nil {
		return nil, errors.New("invalid json format")
	}

	validate := validator.New()

	validate.RegisterValidation("slot", utils.SlotValidater)

	if err := validate.Struct(pdfDownloadRequest); err != nil {
		return nil, errors.New("invalid request format")
	}

	dates, err := utils.GetMiddleDates(pdfDownloadRequest.StartDate, pdfDownloadRequest.EndDate)

	if err != nil {
		return nil, errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	pdf := gopdf.GoPdf{}

	if err := utils.InitPdf(&pdf); err != nil {
		log.Println(err)
		return nil, err
	}

	for _, date := range dates {
		studentAttendanceLogs, err := query.GetStudentsAttendanceLogForPdf(pdfDownloadRequest.UnitId, pdfDownloadRequest.UserId, date, pdfDownloadRequest.Slot)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		if err := utils.GeneratePdf(&pdf, date, strings.ToUpper(pdfDownloadRequest.UnitId), strings.ToUpper(pdfDownloadRequest.Slot), studentAttendanceLogs); err != nil {
			log.Println(err)
			return nil, err
		}

	}

	return &pdf, nil

}

func (repo *studentRepo) DownloadExcel(r *http.Request) (*excelize.File, error) {
	var excelRequest models.ExcelDownloadRequest

	if err := json.NewDecoder(r.Body).Decode(&excelRequest); err != nil {
		return nil, errors.New("invalid json format")
	}

	validate := validator.New()

	validate.RegisterValidation("slot", utils.SlotValidater)

	if err := validate.Struct(excelRequest); err != nil {
		return nil, errors.New("invalid request format")
	}

	dates, err := utils.GetMiddleDates(excelRequest.StartDate, excelRequest.EndDate)

	if err != nil {
		return nil, errors.New("invalid request format")
	}

	days, err := utils.GetMiddleDays(excelRequest.StartDate, excelRequest.EndDate)

	if err != nil {
		return nil, errors.New("invalid request format")
	}

	excelFile := excelize.NewFile()

	sheetName := "attendance"

	_, err = excelFile.NewSheet(sheetName)

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}

	utils.SetExcelHeader(excelFile, sheetName, days)

	query := database.NewQuery(repo.db)

	students, err := query.GetStudentsForExcel(excelRequest.UnitId)

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}

	utils.AddStudentsToExcelFile(excelFile, sheetName, students)

	for col, date := range dates {
		studentsAttendanceStatus, err := query.GetStudentsAttendanceStatusForExcel(excelRequest.UnitId, excelRequest.UserId, date, excelRequest.Slot)

		if err != nil {
			log.Println(err)
			return nil, errors.New("internal server error")
		}

		if err := utils.AddStudentAttendanceStatusToExcelFile(excelFile, sheetName, col+3, studentsAttendanceStatus); err != nil {
			log.Println(err)
			return nil, errors.New("internal server error")
		}
	}

	return excelFile, nil

}
