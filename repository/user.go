package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/database"
	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const OTP_EXPIRE_TIME = 5

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (repo *userRepo) CreateUser(r *http.Request) (string, error) {
	var createUserRequest models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		return "", errors.New("invalid json format")
	}

	validate := validator.New()

	validate.RegisterValidation("strongPassword", utils.StrongPasswordValidator)

	if err := validate.Struct(createUserRequest); err != nil {
		return "", errors.New("invalid request format")
	}

	var user models.User

	hashedPassword, err := utils.HashPassword(createUserRequest.Password)

	if err != nil {
		log.Println(err)
		return "", err
	}

	user.UserId = uuid.NewString()
	user.UserName = createUserRequest.UserName
	user.Password = hashedPassword
	user.Email = createUserRequest.Email

	query := database.NewQuery(repo.db)

	if err := query.CreateUser(&user); err != nil {
		log.Println(err)
		return "", errors.New("internal server error")
	}

	token, err := utils.GenerateToken(user.UserId, user.UserName)

	if err != nil {
		log.Println(err)
		return "", errors.New("internal server error")
	}

	return token, nil

}

func (repo *userRepo) GiveUserAccess(r *http.Request) (bool, error) {
	var giveUserAccessRequest models.GiveUserAccessRequest

	if err := json.NewDecoder(r.Body).Decode(&giveUserAccessRequest); err != nil {
		return false, errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(giveUserAccessRequest); err != nil {
		return false, errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	userName, password, email, err := query.GiveUserAccess(giveUserAccessRequest.UserId)

	if err != nil {
		log.Println(err)
		return false, errors.New("internal server error occurred")
	}

	if err := utils.CheckPassword(password, giveUserAccessRequest.Password); err != nil {
		return false, nil
	}

	if err := utils.SendUserCredentialsToEmail(userName, giveUserAccessRequest.Password, email); err != nil {
		log.Println(err)
		return false, errors.New("internal server error occurred")
	}

	return true, nil

}

func (repo *userRepo) UserLogin(r *http.Request) (bool, string, error) {
	var userLoginRequest models.UserLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&userLoginRequest); err != nil {
		return false, "", errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(userLoginRequest); err != nil {
		return false, "", errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	userId, password, err := query.UserLogin(userLoginRequest.UserName)

	if err != nil {
		log.Println(err)
		return false, "", err
	}

	if err := utils.CheckPassword(password, userLoginRequest.Password); err != nil {
		return false, "", nil
	}

	token, err := utils.GenerateToken(userId, userLoginRequest.UserName)

	if err != nil {
		log.Println(err)
		return false, "", errors.New("internal server error")
	}

	return true, token, nil
}

func (repo *userRepo) GetAllUsers(r *http.Request) ([]*models.User, error) {
	query := database.NewQuery(repo.db)

	users, err := query.GetAllUsers()

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error occurred")
	}

	return users, nil
}

func (repo *userRepo) UpdateNewPassword(r *http.Request) (bool, error) {
	var newPasswordUpdateRequest models.PasswordUpdateRequest
	json.NewDecoder(r.Body).Decode(&newPasswordUpdateRequest)

	validate := validator.New()

	_ = validate.RegisterValidation("strongPassword", utils.StrongPasswordValidator)

	if err := validate.Struct(newPasswordUpdateRequest); err != nil {
		log.Println(err)
		return false, errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)

	isUserIdExists, err := query.CheckUserIdExists(newPasswordUpdateRequest.UserId)

	if err != nil {
		log.Println(err)
		return false, err
	}

	if !isUserIdExists {
		return false, nil
	}

	password, err := utils.HashPassword(newPasswordUpdateRequest.NewPassword)

	if err != nil {
		log.Println(err)
		return false, errors.New("internal server error")
	}

	if err := query.UpdateNewPassword(newPasswordUpdateRequest.UserId, password); err != nil {
		log.Println(err)
		return false, errors.New("internal server error")
	}

	return true, nil
}

func (repo *userRepo) ForgotPassword(r *http.Request) (bool, string, error) {

	var forgotPasswordRequest models.ForgotPasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&forgotPasswordRequest); err != nil {
		return false, "", errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(forgotPasswordRequest); err != nil {
		return false, "", errors.New("invalid request format")
	}
	query := database.NewQuery(repo.db)

	isEmailExists, err := query.CheckUserEmailExists(forgotPasswordRequest.Email)

	if err != nil {
		log.Println(err)
		return false, "", err
	}

	if !isEmailExists {
		return false, "", nil
	}

	otp, err := utils.GenerateOtp()

	if err != nil {
		log.Println(err)
		return false, "", errors.New("internal server error")
	}

	if err := query.StoreOtp(forgotPasswordRequest.Email, otp); err != nil {
		log.Println(err)
		return false, "", errors.New("internal server error")
	}

	go func() {
		time.Sleep(time.Minute * OTP_EXPIRE_TIME)
		if err := query.ClearOtp(forgotPasswordRequest.Email, otp); err != nil {
			log.Println(err)
		}
	}()

	if err := utils.SendOtpToEmail(forgotPasswordRequest.Email, otp, strconv.Itoa(OTP_EXPIRE_TIME)); err != nil {
		log.Println(err)
		return false, "", errors.New("internal server error")
	}
	return true, strconv.Itoa(OTP_EXPIRE_TIME), nil
}

func (repo *userRepo) ValidateOtp(r *http.Request) (string, error) {
	var otpValidationRequest models.ValidateOtpRequest

	if err := json.NewDecoder(r.Body).Decode(&otpValidationRequest); err != nil {
		return "", errors.New("invalid json format")
	}

	validate := validator.New()

	if err := validate.Struct(otpValidationRequest); err != nil {
		return "", errors.New("invalid request format, validation failed!!")
	}

	query := database.NewQuery(repo.db)

	isOtpValid, userId, err := query.IsOtpValid(otpValidationRequest.Email, otpValidationRequest.Otp)

	if err != nil {
		log.Println(err)
		return "", errors.New("internal server error")
	}

	if !isOtpValid {
		return "", errors.New("invalid otp")
	}

	return userId, nil

}

func (repo *userRepo) UpdateTime(r *http.Request) error {
	var updateTimeRequest models.UpdateTimeRequest

	if err := json.NewDecoder(r.Body).Decode(&updateTimeRequest); err != nil {
		return errors.New("invalid json format")
	}

	validate := validator.New()
	validate.RegisterValidation("hhmm", func(fl validator.FieldLevel) bool {
		_, err := time.Parse("15:04", fl.Field().String())
		return err == nil
	})

	if err := validate.Struct(updateTimeRequest); err != nil {
		return errors.New("invalid request format")
	}

	if err := utils.UserUpdateTimeValidater(
		updateTimeRequest.MorningStartTime,
		updateTimeRequest.MorningEndTime,
		updateTimeRequest.AfterNoonStartTime,
		updateTimeRequest.AfterNoonEndTime,
		updateTimeRequest.EveningStartTime,
		updateTimeRequest.EveningEndTime,
	); err != nil {
		log.Println(err)
		return errors.New("invalid request format")
	}

	query := database.NewQuery(repo.db)
	if err := query.UpdateTime(
		updateTimeRequest.UserId,
		updateTimeRequest.MorningStartTime,
		updateTimeRequest.MorningEndTime,
		updateTimeRequest.AfterNoonStartTime,
		updateTimeRequest.AfterNoonEndTime,
		updateTimeRequest.EveningStartTime,
		updateTimeRequest.EveningEndTime,
	); err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}

	return nil
}

func (repo *userRepo) GetBiometricDevicesForRegisterForm(r *http.Request) ([]string, error) {
	userId := mux.Vars(r)["user_id"]

	query := database.NewQuery(repo.db)

	exists, err := query.CheckUserIdExists(userId)

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}

	if !exists {
		return nil, errors.New("user not found")
	}

	units, err := query.GetBiometricDevicesForRegisterForm(userId)

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}
	return units, nil
}

func (repo *userRepo) GetStudentUnitIdsForRegisterForm(r *http.Request) ([]string, error) {
	unitId := mux.Vars(r)["unit_id"]

	query := database.NewQuery(repo.db)

	unitId = strings.ToLower(unitId)

	exists, err := query.CheckBiometricDeviceExists(unitId)

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}

	if !exists {
		return nil, errors.New("biometric device not found")
	}

	studentUnitIds, err := query.GetStudentUnitIdsForRegisterForm(unitId)

	if err != nil {
		log.Println(err)
		return nil, errors.New("internal server error")
	}

	return studentUnitIds, nil
}

func (repo *userRepo) DeleteUser(r *http.Request) error {
	vars := mux.Vars(r)
	userId := vars["user_id"]

	query := database.NewQuery(repo.db)

	exists, err := query.CheckUserIdExists(userId)
	if err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}

	if !exists {
		return errors.New("user not found")
	}

	if err := query.DeleteUser(userId); err != nil {
		log.Println(err)
		return errors.New("internal server error")
	}

	return nil
}
