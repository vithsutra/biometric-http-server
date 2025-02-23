package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type adminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *adminRepo {
	return &adminRepo{
		db,
	}
}

func (repo *adminRepo) CreateAdmin(r *http.Request) (bool, error) {
	var adminRegisterRequest models.AdminRegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&adminRegisterRequest); err != nil {
		return false, errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(adminRegisterRequest); err != nil {
		return false, errors.New("invalid input format")
	}

	rootPassword := os.Getenv("ROOT_PASSWORD")

	if rootPassword == "" {
		log.Println("ROOT_PASSWORD env is empty or missing")
		return false, errors.New("internal server error")
	}

	if adminRegisterRequest.RootPassword != rootPassword {
		return false, nil
	}

	var admin models.Admin

	admin.UserId = uuid.NewString()
	admin.UserName = adminRegisterRequest.UserName

	hashedPassword, err := utils.HashPassword(adminRegisterRequest.Password)

	if err != nil {
		log.Println(err)
		return false, errors.New("internal server error")
	}

	admin.Password = hashedPassword

	query := database.NewQuery(repo.db)

	if err := query.CreateAdmin(&admin); err != nil {
		log.Println(err)
		return false, errors.New("internal server error")
	}

	return true, nil
}

func (repo *adminRepo) AdminLogin(r *http.Request) (bool, error) {
	var adminLoginRequest models.AdminLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&adminLoginRequest); err != nil {
		return false, errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(adminLoginRequest); err != nil {
		return false, errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	password, err := query.GetAdminPassword(adminLoginRequest.UserName)

	if err != nil {
		log.Println(err)
		return false, err
	}

	if err := utils.CheckPassword(password, adminLoginRequest.Password); err != nil {
		return false, nil
	}

	return true, nil

}
