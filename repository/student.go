package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

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

	student.StudentId = uuid.NewString()
	student.StudentUnitId = createStudentRequest.StudentUnitId
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

	if err := query.UpdateStudent(studentUpdateRequest.UnitId, studentUpdateRequest.StudentId, studentUpdateRequest.StudentName, studentUpdateRequest.StudentUsn, studentUpdateRequest.Department); err != nil {
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

	if err := query.DeleteStudent(deleteStudentRequest.UnitId, deleteStudentRequest.StudentId, deleteStudentRequest.StudentUnitId); err != nil {
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

func (repo *studentRepo) DownloadPdf(r *http.Request) error {
	var pdfDownloadRequest models.PdfDownloadRequest

	if err := json.NewDecoder(r.Body).Decode(&pdfDownloadRequest); err != nil {
		return errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(pdfDownloadRequest); err != nil {
		return errors.New("invalid request format")
	}

	midDates, err := utils.GetMiddleDates(pdfDownloadRequest.StartDate, pdfDownloadRequest.EndDate)

	log.Println(err)

	if err != nil {
		return errors.New("invalid request format")
	}

	for _, date := range midDates {
		log.Println(date)
	}

	return nil
}
