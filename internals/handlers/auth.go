package handlers

import (
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
)

type AuthHandler struct {
	authRepo models.AuthInterface
}

func NewAuthHandler(authRepo models.AuthInterface) *AuthHandler {
	return &AuthHandler{
		authRepo,
	}
}

func (ah *AuthHandler) RegisterHandler(w http.ResponseWriter , r *http.Request) {
	
}