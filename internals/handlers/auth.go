package handlers

import (
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
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
	token , err := ah.authRepo.Register(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message":err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string]string{"token":token})
}

func (ah *AuthHandler) LoginHandler(w http.ResponseWriter , r *http.Request) {
	token , err := ah.authRepo.Login(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message":err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string]string{"token":token})
}