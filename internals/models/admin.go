package models

import "net/http"

type AdminRegisterRequest struct {
	RootPassword string `json:"root_password" validate:"required"`
	UserName     string `json:"user_name" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

type Admin struct {
	UserId   string
	UserName string
	Password string
}

type AdminLoginRequest struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminInterface interface {
	CreateAdmin(r *http.Request) (bool, error)
	AdminLogin(r *http.Request) (bool, error)
}
