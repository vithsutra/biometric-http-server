package utils

import (
	"regexp"

	"github.com/go-playground/validator"
)

func TimeValidator(fl validator.FieldLevel) bool {
	time := fl.Field().String()
	return regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`).MatchString(time)
}
