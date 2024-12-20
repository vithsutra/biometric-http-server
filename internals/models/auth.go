package models

import "net/http"

type Auth struct{
	UserId string `json:"user_id,omitempty"`
	Name string `json:"user_name,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthInterface interface{
	Register(r *http.Request) (string,error)
	Login(r *http.Request) (string,error)
}