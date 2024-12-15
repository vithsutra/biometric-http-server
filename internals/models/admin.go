package models

import "net/http"

type Admin struct{
	UserId string `json:"user_id"`
	UserName string `json:"user_name"`
	Email string `json:"-"`
	Password string `json:"-"`
}

type AdminInterface interface{
	FetchAllUsers(*http.Request) ([]Admin , error)
	GiveUserAccess(*http.Request) error
}