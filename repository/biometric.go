package repository

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
	"github.com/gorilla/mux"
)


type BiometricRepo struct {
	db *sql.DB
}

func NewBiometricRepo(db *sql.DB) *BiometricRepo {
	return &BiometricRepo{
		db,
	}
}

func (br *BiometricRepo) NewBiometricDevice(r *http.Request) error {
	var newDevice models.Biometric
	if err := utils.Decode(r , &newDevice); err != nil {
		return err
	}
	if newDevice.UnitId == "" || newDevice.UserId == "" {
		return fmt.Errorf("enter a valid user")
	}
	query := database.NewQuery(br.db)
	if err := query.NewBiometricDevice(newDevice); err != nil {
		return err
	}
	return nil
}

func (br *BiometricRepo) FetchAllBiometrics(r *http.Request) ([]models.Biometric , error) {
	var userId = mux.Vars(r)["userid"]
	query := database.NewQuery(br.db)
	data , err := query.FetchAllBiometrics(userId)
	if err != nil {
		return nil,err
	}
	return data , nil
}

func (br *BiometricRepo) DeleteBiometricMachine(r *http.Request) error {
	var unitId  = mux.Vars(r)["unitid"]
	query := database.NewQuery(br.db)
	if err := query.DeleteBiometricMachine(unitId); err != nil {
		return err
	}
	return nil
}