package utils

import (
	"encoding/json"
	"net/http"
)

func Encode(w http.ResponseWriter , input any) error {
	if err := json.NewEncoder(w).Encode(input); err != nil {
		return err
	}
	return nil
}

func Decode(r *http.Request , input any) error {
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		return err
	}
	return nil
}