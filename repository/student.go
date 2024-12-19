package repository

import (
	"bytes"
	"database/sql"
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type StudentRepo struct {
	db *sql.DB
}

func NewStudentRepo(db *sql.DB) *StudentRepo {
	return &StudentRepo{
		db,
	}
}

func (sr *StudentRepo) GenerateStudentAttendenceReport(r *http.Request) (*bytes.Buffer, error) {
	var details models.Student
	if err := utils.Decode(r , &details); err != nil {
		return nil,err
	}
	query := database.NewQuery(sr.db)
	data ,err := query.GenerateStudentAttendenceReport(details.UnitId , details.Date)
	if err != nil {
		return nil,err
	}
	buf , err := utils.GeneratePDF(data)
	if err != nil {
		return nil , err
	}
	return buf , nil
}

func(sr *StudentRepo) NewStudent(r *http.Request) error {
	var newStudent models.Student
	if err := utils.Decode(r , &newStudent); err != nil {
		return err
	}
	newStudent.StudentId = uuid.NewString()
	query := database.NewQuery(sr.db)
	if err := query.NewStudent(newStudent); err != nil {
		return err
	}
	return nil
}

func (sr *StudentRepo) DeleteStudent(r *http.Request) error {
	var studentId string = mux.Vars(r)["studentid"]
	var unitId string = mux.Vars(r)["unitid"]
	var studentUnitID string = mux.Vars(r)["studentunitid"]
	query := database.NewQuery(sr.db)
	err := query.DeleteStudent(unitId , studentUnitID , studentId)
	if err != nil {
		return err
	}
	return nil
}

func (sr *StudentRepo) UpdateStudent(r *http.Request) error {
	var student models.Student
	if err := utils.Decode(r , student); err != nil {
		return err
	}
	query := database.NewQuery(sr.db)
	if err := query.UpdateStudent(student); err != nil {
		return err
	}
	return nil
}

func (sr *StudentRepo) FetchStudentDetails(r *http.Request) ([]models.Student , error) {
	var unitId = mux.Vars(r)["unitid"]
	query := database.NewQuery(sr.db)
	data , err := query.FetchStudentDetails(unitId)
	if err != nil {
		return nil,err
	}
	return data, nil
}

func (sr *StudentRepo) FetchStudentLogHistory(r *http.Request) ([]models.Student , error) {
	var studentId = mux.Vars(r)["studentid"]
	query := database.NewQuery(sr.db)
	data , err := query.FetchStudentLogHistory(studentId)
	if err != nil {
		return nil , err
	}
	return data , nil
}