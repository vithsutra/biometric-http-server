package utils

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator"
)

func StrongPasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(string(ch)):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial

}
