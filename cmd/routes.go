package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InitilizeHttpRouters() http.Handler {
	router := mux.NewRouter()

	return router
}