package handlers

import (
	"fmt"
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
)


type ExcelHandler struct{
	excelRepo models.ExcelInterface
}

func NewExcelHandler(excelRepo models.ExcelInterface) *ExcelHandler {
	return &ExcelHandler{
		excelRepo,
	}
}

func (eh *ExcelHandler) GenerateExcelReportHandler(w http.ResponseWriter , r *http.Request) {
	file , name ,err := eh.excelRepo.GenerateExcelReport(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message" : err.Error()})
		return
	}
	sheetName := fmt.Sprintf("attachment;filename=%s" , name)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", sheetName)
	w.Header().Set("File-Name", name)
	if err := file.Write(w); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w , map[string]string{"message" : err.Error()})
		return
	}
}