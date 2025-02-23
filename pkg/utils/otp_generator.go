package utils

import (
	"crypto/rand"
	"math/big"
)

const digits = "0123456789"
const otpLength = 6

func GenerateOtp() (string, error) {
	otp := ""

	for i := 0; i < otpLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))

		if err != nil {
			return "", err
		}

		otp += string(digits[n.Int64()])
	}

	return otp, nil
}
