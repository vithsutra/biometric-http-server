package repository

import (
	"bytes"
	"database/sql"
	"net/http"
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
	
	return nil , nil
}