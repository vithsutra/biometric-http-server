package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

	var student models.Student

	unitId := strings.ToLower(createStudentRequest.UnitId)
	student.StudentId = uuid.NewString()
	student.StudentUnitId = unitId
	student.StudentName = createStudentRequest.StudentName
	student.StudentUsn = createStudentRequest.StudentUsn
	student.Department = createStudentRequest.Department

	query := database.NewQuery(repo.db)

	if err := query.CreateNewStudent(&student, createStudentRequest.UnitId, createStudentRequest.FingerprintData); err != nil {
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

func (repo *studentRepo) DownloadPdf(r *http.Request) ([][]*models.PdfFormat, error) {
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

	log.Println(dates)

	return nil, nil

	query := database.NewQuery(repo.db)

	logs, err := query.GetStudentsAttendanceLogForPdf("vs242s38", "2f706016-10b4-406c-ba00-fb6fa2cb1374", "2025-02-26", "full")

	if err != nil {
		log.Fatal(err)
	}

	for _, l := range logs {
		fmt.Println(l)
	}
	return nil, nil

}
