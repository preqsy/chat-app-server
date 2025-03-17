package core

import (
	models "chat_app_server/model"
	"context"
	"errors"
	"fmt"
)

func (s *Service) SendFriendRequest(ctx context.Context, sender *models.AuthUser, receiverId uint) (*models.AuthUser, error) {

	receiver, err := s.datastore.GetUserById(ctx, receiverId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	fmt.Println("receiver", receiver)
	_, err = s.neo4jService.SendFriendRequest(ctx, sender, receiver)

	if err != nil {
		return nil, err
	}
	return receiver, nil
}
