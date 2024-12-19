package repository

import (
	"bytes"
	"database/sql"
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
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
	if err := utils.Decode(r , details); err != nil {
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
	if err := utils.Decode(r , newStudent); err != nil {
		return err
	}
	query := database.NewQuery(sr.db)
	query.
}