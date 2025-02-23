package utils

import (
	"regexp"

	"github.com/go-playground/validator"
)

func TimeValidator(fl validator.FieldLevel) bool {
	time := fl.Field().String()
	return regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,9})? \+\d{4} UTC$`).MatchString(time)
}
