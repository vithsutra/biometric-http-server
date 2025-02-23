package models

import "net/http"

type PasswordUpdateRequest struct {
	UserId      string `json:"user_id" validate:"required"`
	NewPassword string `json:"password" validate:"required,strongPassword"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ValidateOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required"`
}

type CreateUserRequest struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required,strongPassword"`
	Email    string `json:"email" validate:"required,email"`
}

type UserLoginRequest struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateTimeRequest struct {
	MorningStartTime   string `json:"morning_start_time" validate:"required,utcTimeFormat"`
	MorningEndTime     string `json:"morning_end_time" validate:"required,utcTimeFormat"`
	AfterNoonStartTime string `json:"afternoon_start_time" validate:"required,utcTimeFormat"`
	AfterNoonEndTime   string `json:"afternoon_end_time" validate:"required,utcTimeFormat"`
	EveningStartTime   string `json:"evening_start_time" validate:"required,utcTimeFormat"`
	EveningEndTime     string `json:"evening_end_time" validate:"required,utcTimeFormat"`
}

type User struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type GiveUserAccessRequest struct {
	UserId   string `json:"user_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserInterface interface {
	CreateUser(r *http.Request) error
	GiveUserAccess(r *http.Request) (bool, error)
	UserLogin(r *http.Request) (bool, string, error)
	GetAllUsers(r *http.Request) ([]*User, error)
	UpdateNewPassword(*http.Request) error
	ForgotPassword(r *http.Request) (string, error)
	ValidateOtp(r *http.Request) error
	UpdateTime(r *http.Request) error
	GetBiometricDevicesForRegisterForm(r *http.Request) ([]string, error)
	GetStudentUnitIdsForRegisterForm(r *http.Request) ([]string, error)
}
