package jwt_utils

import (
	"chat_app_server/config"
	datastore "chat_app_server/database"
	models "chat_app_server/model"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type JWTUtils struct {
	datastore datastore.Datastore
	logger    *logrus.Logger
}

func InitializeJWTUtils(datastore datastore.Datastore, logger *logrus.Logger) *JWTUtils {
	return &JWTUtils{datastore: datastore, logger: logger}
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

func (j *JWTUtils) VerifyAccessToken(tokenString string) (uint, error) {
	secret := config.GetSecrets()

	secretKey := []byte(secret.JwtSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		j.logger.Error("Invalid Token")
		return 0, err
	}
	userID := token.Claims.(jwt.MapClaims)["user_id"].(float64)
	if userID == 0 {
		j.logger.Errorf("Invalid User")
		return 0, err
	}

	if !token.Valid {
		j.logger.Error("Token Expired")
		return 0, err
	}

	return uint(userID), nil
}

func (j *JWTUtils) GetCurrentAuthUser(ctx context.Context) (*models.AuthUser, error) {
	request, ok := ctx.Value("request").(*http.Request)
	if !ok {
		return nil, fmt.Errorf("request not found")
	}
	token := request.Header.Get("authorization")
	cleanToken := strings.TrimPrefix(token, "Bearer ")

	userId, err := j.VerifyAccessToken(cleanToken)
	if err != nil {
		return nil, err
	}
	data, err := j.datastore.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return data, nil
}
