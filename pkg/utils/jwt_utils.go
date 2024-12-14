package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)


func GenerateToken(userId string , userName string) (string , error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        userId,
		"user_name": userName,
		"expiry":    time.Now().Add(365 * 24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("vsense"))
	if err != nil {
		return "", err
	}
	return tokenString,nil
}