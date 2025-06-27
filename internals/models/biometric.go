package models

import "net/http"

type CreateBiometricRequest struct {
	UnitId string `json:"unit_id" validate:"required"`
	UserId string `json:"user_id" validate:"required"`
	Label  string `json:"label" validate:"required"`
}

type ClearBiometricDataRequest struct {
	UserId string `json:"user_id" validate:"required"`
	UnitId string `json:"unit_id" validate:"required"`
}

type UpdateBiometricLabelRequest struct {
	UnitId string `json:"unit_id" validate:"required"`
	Label  string `json:"label" validate:"required"`
}

type Biometric struct {
	UnitId string `json:"unit_id"`
	UserId string `json:"user_id"`
	Online bool   `json:"online"`
	Label  string `json:"label"`
}

type BiometricInterface interface {
	CreateBiometricDevice(r *http.Request) error
	GetBiometricDevices(r *http.Request) ([]*Biometric, error)
	UpdateBiometricLabel(r *http.Request) error
	DeleteBiometricDevice(r *http.Request) error
}
