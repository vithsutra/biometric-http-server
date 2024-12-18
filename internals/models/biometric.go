package models

import "net/http"

type Biometric struct{
	UnitId string `json:"unit_id"`
	UserId string `json:"user_id"`
	Status bool `json:"online"`
}

type BiometricInterface interface {
	FetchAllBiometrics(*http.Request) ([]Biometric , error)
	DeleteBiometricMachine(*http.Request) error
	NewBiometricDevice(*http.Request) error
}
