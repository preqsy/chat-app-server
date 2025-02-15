package auth

import (
	datastore "chat_app_server/database"
	"chat_app_server/jwt_utils"
	models "chat_app_server/model"
	"chat_app_server/utils"
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
	user.Password, _ = utils.HashPassword(user.Password)
	savedUser, err := s.datastore.SaveUser(user)
	if err != nil {
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

func (s *Service) LoginUser(payload *models.AuthUserLogin) (string, error) {
	user, err := s.datastore.GetUserByEmail(payload.Email)
	if err != nil {
		return "", err
	}
	err = utils.VerifyPassword(user.Password, payload.Password)

	if err != nil {
		return "", err
	}
	token, err := jwt_utils.GenerateAccessToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil

}
