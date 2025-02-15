package jwt_utils

import (
	"chat_app_server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
