package auth

import (
	datastore "chat_app_server/database"
	models "chat_app_server/model"
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

	user.Password = "password"
	savedUser, err := s.datastore.SaveUser(user)
	if err != nil {
		return nil, err
	}
	return savedUser, nil
}
