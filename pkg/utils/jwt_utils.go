package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(userId string, userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        userId,
		"user_name": userName,
		"exp":       time.Now().Add(365 * 24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("vsense"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateAdminToken(userID string, userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id":   userID,
		"admin_name": userName,
		"exp":        time.Now().Add(365 * 24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("vsense"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("vsense"), nil
	})
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return err
	}

	if claims.Valid() != nil {
		return fmt.Errorf("token is invalid")
	}

	if claims["user_id"] == "" || claims["user_name"] == "" {
		return fmt.Errorf("user_id or user_name is invalid")
	}

	return nil
}

func ValidateAdminToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("vsense"), nil
	})
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return err
	}

	if claims.Valid() != nil {
		return fmt.Errorf("token is invalid")
	}

	if claims["admin_id"] == "" || claims["admin_name"] == "" {
		return fmt.Errorf("admin_id or admin_name is invalid")
	}

	return nil
}
