package auth

import (
	datastore "chat_app_server/database"
	"chat_app_server/jwt_utils"
	models "chat_app_server/model"
	"chat_app_server/utils"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Service struct {
	datastore datastore.Datastore
}

func CoreService(datastore datastore.Datastore) *Service {
	return &Service{
		datastore: datastore,
	}
}

func (s *Service) SaveUser(user *models.AuthUser) (*models.AuthUserRegisterResponse, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	_, err := s.datastore.GetUserByEmail(user.Email)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user.Password, _ = utils.HashPassword(user.Password)
			user.Email = strings.ToLower(user.Email)
			savedUser, err := s.datastore.SaveUser(user)
			if err != nil {
				if strings.Contains(err.Error(), "uni_auth_users_email") {
					return nil, fmt.Errorf("account with email %s already exists", user.Email)
				}
				return nil, err
			}
			token, err := jwt_utils.GenerateAccessToken(savedUser.ID)
			if err != nil {
				return nil, err
			}
			response := &models.AuthUserRegisterResponse{
				AuthUser: *savedUser,
				Token:    token,
			}
			return response, nil
		}
		return nil, err
	}

	return nil, fmt.Errorf("account with email %s already exists", user.Email)
}

func (s *Service) LoginUser(payload *models.AuthUserLogin) (string, error) {
	user, err := s.datastore.GetUserByEmail(payload.Email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	err = utils.VerifyPassword(user.Password, payload.Password)

	if err != nil {
		return "Invalid Credentials", fmt.Errorf("invalid credentials")
	}
	token, err := jwt_utils.GenerateAccessToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (s *Service) GetCurrentUser(email string) (*models.AuthUser, error) {
	user, err := s.datastore.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return user, nil
}
