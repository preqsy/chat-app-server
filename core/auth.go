package auth

import (
	datastore "chat_app_server/database"
	models "chat_app_server/model"
	"fmt"
	// "github.com/google/uuid"
)

type Service struct {
	datastore datastore.Datastore
}

func CoreService(datastore datastore.Datastore) *Service {
	fmt.Printf("CoreService: Storing datastore reference at memory address %p\n", &datastore)
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

	fmt.Printf("SaveUser: Created new user struct at memory address %p\n", &newUser)
	fmt.Printf("SaveUser: Passing newUser to datastore at memory address %p\n", &newUser)

	err := s.datastore.SaveUser(&newUser)
	if err != nil {
		return err
	}
	return nil
}
