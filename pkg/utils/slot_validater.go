package utils

import "github.com/go-playground/validator"

func SlotValidater(fl validator.FieldLevel) bool {
	slot := fl.Field().String()

	if slot == "morning" || slot == "afternoon" || slot == "full" {
		return true
	}
	return false
}
