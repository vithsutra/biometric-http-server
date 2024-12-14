package repository

import (
	"database/sql"
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db,
	}
}

func (ar *AuthRepo) Register(r *http.Request) (string , error) {
	var vars = mux.Vars(r)["id"]
	var newUser models.Auth

	if err := utils.Decode(r , &newUser); err != nil {
		return "" , err
	}

	newUser.UserId = uuid.NewString()
	hash , err := utils.HashPassword(newUser.Password)
	if err != nil {
		return "" , err
	}
	newUser.Password = hash

	token , err := utils.GenerateToken(newUser.UserId , newUser.Name)
	if err != nil {
		return "" , err
	}

	query := database.NewQuery(ar.db)
	if err := query.Register(newUser , vars); err != nil {
		return "",err
	}

	return token,nil
}

func (ar *AuthRepo) Login(r *http.Request) (string , error) {
	var vars = mux.Vars(r)["id"]
	var reqUser models.Auth

	if err := utils.Decode(r , &reqUser); err != nil {
		return "" , err
	}

	query := database.NewQuery(ar.db)
	dbUser  , err := query.Login(reqUser , vars)
	if err != nil {
		return "",err
	}
	if err := utils.CheckPassword(dbUser.Password , reqUser.Password); err != nil {
		return "", err
	}
	token , err := utils.GenerateToken(dbUser.UserId , dbUser.Name)
	if err != nil {
		return "",err
	}
	return token, nil
}