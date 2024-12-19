package handlers

import (
	"fmt"
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
)

type StudentHandler struct {
	studentRepo models.StudentInterface
}

func NewStudentHandler(studentRepo models.StudentInterface) *StudentHandler {
	return &StudentHandler{
		studentRepo,
	}
}


func (sh *StudentHandler) GenerateStudentAttendenceReportHandler(w http.ResponseWriter, r *http.Request) {
    buff, err := sh.studentRepo.GenerateStudentAttendenceReport(r)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        utils.Encode(w, map[string]string{"message": err.Error()})
        return
    }

    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "inline; filename=attendance_report.pdf")
    w.Header().Set("Content-Length", fmt.Sprintf("%d", buff.Len()))

    if _, err := w.Write(buff.Bytes()); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        utils.Encode(w, map[string]string{"message": "Failed to send PDF"})
        return
    }
}

func (sh *StudentHandler) DeleteStudentHandler(w http.ResponseWriter , r *http.Request) {
	if err := sh.studentRepo.DeleteStudent(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message":err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string]string{"message":"Success"})
}

func (sh *StudentHandler) UpdateStudentHandler(w http.ResponseWriter , r *http.Request) {
	if err := sh.studentRepo.UpdateStudent(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message":err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string]string{"message":"Success"})
}

func (sh *StudentHandler) NewStudentHandler(w http.ResponseWriter , r *http.Request) {
	if err := sh.studentRepo.NewStudent(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message":err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string]string{"message":"Success"})
}

func (sh *StudentHandler) FetchStudentDetails(w http.ResponseWriter , r *http.Request) {
	data , err := sh.studentRepo.FetchStudentDetails(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message":err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string][]models.Student{"students":data})
}

func (sh *StudentHandler) FetchStudentLogHistoryHandler(w http.ResponseWriter , r *http.Request) {
	data , err := sh.studentRepo.FetchStudentLogHistory(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message":err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w , map[string][]models.Student{"logs":data})
}


