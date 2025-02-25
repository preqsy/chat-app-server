package jwt_utils

import (
	"chat_app_server/config"
	datastore "chat_app_server/database"
	models "chat_app_server/model"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type JWTUtils struct {
	db datastore.Datastore
}

func InitDB(database datastore.Datastore) *JWTUtils {
	return &JWTUtils{db: database}
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

func VerifyAccessToken(tokenString string) (uint, error) {
	secret := config.GetSecrets()

	secretKey := []byte(secret.JwtSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		logrus.Error("Invalid Token")
		return 0, err
	}
	userID := token.Claims.(jwt.MapClaims)["user_id"].(float64)
	if userID == 0 {
		return 0, fmt.Errorf("Invalid User")
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return uint(userID), nil
}

func (j *JWTUtils) GetCurrentAuthUser(token string) (*models.AuthUser, error) {
	userId, err := VerifyAccessToken(token)
	if err != nil {
		return nil, err
	}
	data, err := j.db.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return data, nil
}
