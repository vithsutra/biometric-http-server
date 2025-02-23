package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
)

type biometricHandler struct {
	repo models.BiometricInterface
}

func NewBiometricHandler(repo models.BiometricInterface) *biometricHandler {
	return &biometricHandler{
		repo,
	}
}

func (h *biometricHandler) CreateBiometricDeviceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := h.repo.CreateBiometricDevice(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "biometric device created successfully"})
}

func (h *biometricHandler) GetBiometricDevicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	units, err := h.repo.GetBiometricDevices(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if units == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"units": []interface{}{},
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]*models.Biometric{
		"units": units,
	})
}

func (h *biometricHandler) UpdateBiometricLabelHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := h.repo.UpdateBiometricLabel(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "biometric device label updated successfully"})
}

func (h *biometricHandler) DeleteBiometricDeviceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := h.repo.DeleteBiometricDevice(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "biometric device deleted successfully"})
}

func (h *biometricHandler) ClearBiometricDeviceDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := h.repo.ClearBiometricDeviceData(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "biometric device data cleared successfully"})
}
