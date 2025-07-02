package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	// "strconv"
	"strings"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/signintech/gopdf"
)

type studentRepo struct {
	db *sql.DB
}

func NewStudentRepo(db *sql.DB) *studentRepo {
	return &studentRepo{
		db,
	}
}

func (repo *studentRepo) CreateNewStudent(r *http.Request) error {
	var createStudentRequest models.CreateStudentRequest

	if err := json.NewDecoder(r.Body).Decode(&createStudentRequest); err != nil {
		return errors.New("invalid json format")
	}

	validate := validator.New()
	if err := validate.Struct(createStudentRequest); err != nil {
		return errors.New("invalid request format")
	}

	if len(createStudentRequest.FingerprintData) != 6 {
		return fmt.Errorf("invalid fingerprint data: expected 6, got %d", len(createStudentRequest.FingerprintData))
	}


	query := database.NewQuery(repo.db)

	availableStudentIDs, UnitFound, err := query.GetAvailableStudentUnitIDs(createStudentRequest.UnitId)
	if err != nil {
		return err
	} else if !UnitFound {
		return fmt.Errorf("biometric device not found")
	}

	// isStudentUnitIdExists, err := query.CheckStudentUnitIdExists(createStudentRequest.UnitId, createStudentRequest.StudentUnitId)

	// if err != nil {
	// 	log.Println(err)
	// 	return errors.New("internal server error")
	// }

	// if isStudentUnitIdExists {
	// 	return errors.New("student unit id already exists")
	// }

	var student models.Student

	unitId := strings.ToLower(createStudentRequest.UnitId)

	student.StudentId = uuid.NewString()
	// student.StudentUnitId = createStudentRequest.StudentUnitId
	student.StudentName = createStudentRequest.StudentName
	student.StudentUsn = createStudentRequest.StudentUsn
	student.Department = createStudentRequest.Department

	student.StudentUnitId = make([]string, 6)
	for i := 0; i < 6; i++ {
		student.StudentUnitId[i] = availableStudentIDs[i]
	}

	if err := query.CreateNewStudent(&student, unitId, createStudentRequest.FingerprintData); err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}

	return nil
}

func (repo *studentRepo) UpdateStudentDetails(r *http.Request) error {
	var studentUpdateRequest models.UpdateStudentRequest

	if err := json.NewDecoder(r.Body).Decode(&studentUpdateRequest); err != nil {
		return errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(studentUpdateRequest); err != nil {
		return errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	unitId := strings.ToLower(studentUpdateRequest.UnitId)
	if err := query.UpdateStudent(unitId, studentUpdateRequest.StudentId, studentUpdateRequest.StudentName, studentUpdateRequest.StudentUsn, studentUpdateRequest.Department); err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}

	return nil
}

func (repo *studentRepo) DeleteStudent(r *http.Request) error {
	var deleteStudentRequest models.DeleteStudentRequest

	if err := json.NewDecoder(r.Body).Decode(&deleteStudentRequest); err != nil {
		return errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(deleteStudentRequest); err != nil {
		return errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	unitId := strings.ToLower(deleteStudentRequest.UnitId)

	err := query.DeleteStudent(unitId, deleteStudentRequest.StudentId)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// func (repo *studentRepo) GetStudentDetails(r *http.Request) ([]*models.Student, int, int, int, error) {
func (repo *studentRepo) GetStudentDetails(r *http.Request) ([]*models.Student, error) {
	vars := mux.Vars(r)

	unitId := vars["unit_id"]

	// limit := r.URL.Query().Get("limit")
	// if limit == "" {
	// 	limit = "10"
	// }

	// page := r.URL.Query().Get("page")
	// if page == "" {
	// 	page = "1"
	// }

	// limit_int, err := strconv.Atoi(limit)
	// if err != nil {
	// 	return nil, -1, -1, -1, errors.New("invalid limit")
	// }

	// page_int, err := strconv.Atoi(page)
	// if err != nil {
	// 	return nil, -1, -1, -1, errors.New("invalid page")
	// }

	query := database.NewQuery(repo.db)

	exists, err := query.CheckBiometricDeviceExists(unitId)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}
	if !exists {
		return nil, errors.New("biometric device not found")
	}

	// offset := (page_int - 1) * limit_int

	// students, total_students, err := query.GetStudentDetails(unitId, limit_int, offset)

	students, err := query.GetStudentDetails(unitId)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}

	// return students, total_students, limit_int, page_int, nil
	return students, nil
}

func (repo *studentRepo) GetStudentLogs(r *http.Request) ([]*models.StudentAttendanceLog, error) {
	studentId := mux.Vars(r)["student_id"]

	query := database.NewQuery(repo.db)

	logs, err := query.GetStudentLogs(studentId)

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}

	return logs, nil

}

func (repo *studentRepo) DownloadPdf(r *http.Request) (*gopdf.GoPdf, error) {
	var pdfDownloadRequest models.PdfDownloadRequest

	if err := json.NewDecoder(r.Body).Decode(&pdfDownloadRequest); err != nil {
		return nil, errors.New("invalid json format")
	}

	validate := validator.New()

	validate.RegisterValidation("slot", utils.SlotValidater)

	if err := validate.Struct(pdfDownloadRequest); err != nil {
		return nil, errors.New("invalid request format")
	}

	dates, err := utils.GetMiddleDates(pdfDownloadRequest.StartDate, pdfDownloadRequest.EndDate)

	if err != nil {
		return nil, errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	pdf := gopdf.GoPdf{}

	if err := utils.InitPdf(&pdf); err != nil {
		log.Println(err)
		return nil, err
	}

	userTimeChann := make(chan *models.UserTime)

	go func() {
		userTime, err := query.GetUserStandardTime(pdfDownloadRequest.UserId)
		if err != nil {
			userTimeChann <- nil
			return
		}

		userTimeChann <- userTime
		return
	}()

	studentsCount, err := query.GetStudentsCountFromUnit(pdfDownloadRequest.UnitId)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	pdfFormatsChann := make(chan map[string]*models.PdfFormat)

	go func() {
		pdfFormats, err := query.GetStudentsForPdf(pdfDownloadRequest.UnitId, studentsCount)

		if err != nil {
			pdfFormatsChann <- nil
			return
		}

		pdfFormatsChann <- pdfFormats
		return
	}()

	userTime := <-userTimeChann

	if userTime == nil {
		return nil, errors.New("error occurred with database")
	}

	pdfFormats := <-pdfFormatsChann

	if pdfFormats == nil {
		return nil, errors.New("error occurred with database")
	}

	for _, date := range dates {
		if err := query.GetStudentsAttendanceLogForPdf(studentsCount, userTime, pdfFormats, date, pdfDownloadRequest.Slot); err != nil {
			log.Println(err)
			return nil, err
		}

		if err := utils.GeneratePdf(&pdf, date, strings.ToUpper(pdfDownloadRequest.UnitId), strings.ToUpper(pdfDownloadRequest.Slot), pdfFormats); err != nil {
			log.Println(err)
			return nil, err
		}

	}

	return &pdf, nil

}
