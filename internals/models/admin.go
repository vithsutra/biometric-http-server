package models

import "net/http"

type Admin struct{
	UserId string `json:"user_id"`
	UserName string `json:"user_name"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type AdminInterface interface{
	FetchAllUsers() ([]Admin , error)
	GiveUserAccess(*http.Request) error
}