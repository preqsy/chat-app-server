package auth

import (
	datastore "chat_app_server/database"
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
	token, err := utils.GenerateAccessToken(savedUser.ID)
	if err != nil {
		return nil, err
	}
	response := &models.AuthUserRegisterResponse{
		AuthUser: *savedUser,
		Token:    token,
	}
	return response, nil
}
