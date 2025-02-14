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

func (s *Service) SaveUser(user *models.AuthUser) (*models.AuthUser, error) {

	user.Password, _ = utils.HashPassword(user.Password)
	if err := user.Validate(); err != nil {
		return nil, err
	}
	savedUser, err := s.datastore.SaveUser(user)
	if err != nil {
		return nil, err
	}
	return savedUser, nil
}
