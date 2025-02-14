package utils

import (
	"chat_app_server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashePassword), nil
}

func GenerateAccessToken(userID uint) (string, error) {
	secret := config.GetSecrets()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	secretBytes := []byte(secret.JwtSecret)
	tokenString, err := token.SignedString(secretBytes)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
