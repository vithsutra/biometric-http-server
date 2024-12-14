package models

import "net/http"

type Auth struct{
	UserId string `json:"user_id"`
	Name string `json:"user_name"`
	Password string `json:"password"`
}

type AuthInterface interface{
	Register(r *http.Request) (string,error)
	Login(r *http.Request) (string,error)
}