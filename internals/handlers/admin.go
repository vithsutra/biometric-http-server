package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
)

type adminHandler struct {
	adminInterface models.AdminInterface
}

func NewAdminHandler(adminInterface models.AdminInterface) *adminHandler {
	return &adminHandler{
		adminInterface,
	}
}

func (h *adminHandler) CreateAdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	isRootPasswordCorrect, err := h.adminInterface.CreateAdmin(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if !isRootPasswordCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "incorrect root password"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "admin created successfully"})
}

func (h *adminHandler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	isPasswordCorrect, token, err := h.adminInterface.AdminLogin(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if !isPasswordCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "incorrect password"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "login successfull", "token": token})
}
