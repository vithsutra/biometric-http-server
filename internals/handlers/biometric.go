package handlers

import (
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
)

type BiometricHandler struct {
	biometricRepo models.BiometricInterface
}

func NewBiometricHandler(biometricRepo models.BiometricInterface) *BiometricHandler {
	return &BiometricHandler{
		biometricRepo,
	}
}

func (bh *BiometricHandler) NewBiometricDeviceHandler(w http.ResponseWriter , r *http.Request) {
	err := bh.biometricRepo.NewBiometricDevice(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message" : err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string]string{"message" : "Success..."})
}

func (bh *BiometricHandler) FetchAllBiometricsHandler(w http.ResponseWriter , r *http.Request) {
	data , err := bh.biometricRepo.FetchAllBiometrics(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message" : err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string][]models.Biometric{"devices" : data})
}

func (bh *BiometricHandler) DeleteBiometricMachineHandler(w http.ResponseWriter , r *http.Request) {
	err := bh.biometricRepo.DeleteBiometricMachine(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message" : err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string]string{"message" : "Success..."})
}