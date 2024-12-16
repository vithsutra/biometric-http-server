package handlers

import (
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
)

type AdminHandler struct{
	adminInterface models.AdminInterface
}

func NewAdminHandler(adminInterface models.AdminInterface) *AdminHandler {
	return &AdminHandler{
		adminInterface,
	}
}

func (ah *AdminHandler) FetchAllUsersHandler(w http.ResponseWriter , r *http.Request){
	users , err := ah.adminInterface.FetchAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message":err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string][]models.Admin{"users":users})
}
func (ah *AdminHandler) GiveUserAccessHandler(w http.ResponseWriter , r *http.Request){
	err := ah.adminInterface.GiveUserAccess(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message":err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string]string{"message":"Email sent successfully"})
}