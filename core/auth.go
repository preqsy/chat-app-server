package auth

import (
	datastore "chat_app_server/database"
	models "chat_app_server/model"
	// "github.com/google/uuid"
)

type Service struct {
	datastore datastore.Datastore
}

func CoreService(datastore datastore.Datastore) *Service {
	return &Service{
		datastore: datastore,
	}
}

func (s *Service) SaveUser(user *models.AuthUserCreate) error {
	var newUser models.AuthUser

	newUser.Email = user.Email
	newUser.Password = user.Password
	newUser.FirstName = user.FirstName
	newUser.LastName = user.LastName
	newUser.Username = user.Username

	err := s.datastore.SaveUser(&newUser)
	if err != nil {
		return err
	}
	return nil
}
