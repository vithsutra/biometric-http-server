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
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type biometricRepo struct {
	db *sql.DB
}

func NewBiometricRepo(db *sql.DB) *biometricRepo {
	return &biometricRepo{
		db,
	}
}

func (repo *biometricRepo) CreateBiometricDevice(r *http.Request) error {
	var createBiometricDeviceRequest models.CreateBiometricRequest

	if err := json.NewDecoder(r.Body).Decode(&createBiometricDeviceRequest); err != nil {
		return errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(createBiometricDeviceRequest); err != nil {
		return errors.New("invalid request format")
	}

	var biometric models.Biometric

	unitId := strings.ToLower(createBiometricDeviceRequest.UnitId)

	biometric.UserId = createBiometricDeviceRequest.UserId
	biometric.UnitId = unitId
	biometric.Label = createBiometricDeviceRequest.Label
	biometric.Online = false

	query := database.NewQuery(repo.db)

	if err := query.CreateBiometricDevice(&biometric); err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}

	return nil
}

func (repo *biometricRepo) GetBiometricDevices(r *http.Request) ([]*models.Biometric, error) {
	vars := mux.Vars(r)
	userId := vars["user_id"]

	query := database.NewQuery(repo.db)

	biometrics, err := query.GetBiometricDevices(userId)

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error occurred")
	}

	return biometrics, nil
}

func (repo *biometricRepo) UpdateBiometricLabel(r *http.Request) error {
	var biometricLabelUpdateRequest models.UpdateBiometricLabelRequest

	if err := json.NewDecoder(r.Body).Decode(&biometricLabelUpdateRequest); err != nil {
		return errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(biometricLabelUpdateRequest); err != nil {
		return errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	if err := query.UpdateBiometricLabel(biometricLabelUpdateRequest.UnitId, biometricLabelUpdateRequest.Label); err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}

	return nil
}

func (repo *biometricRepo) DeleteBiometricDevice(r *http.Request) error {
	vars := mux.Vars(r)

	unitId := vars["unit_id"]

	query := database.NewQuery(repo.db)

	if err := query.DeleteBiometricDevice(unitId); err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}
	return nil
}

func (repo *biometricRepo) ClearBiometricDeviceData(r *http.Request) error {
	var clearBiometricRequest models.ClearBiometricDataRequest

	if err := json.NewDecoder(r.Body).Decode(&clearBiometricRequest); err != nil {
		return errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(clearBiometricRequest); err != nil {
		return errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	if err := query.ClearBiometricDeviceData(clearBiometricRequest.UserId, clearBiometricRequest.UnitId); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
