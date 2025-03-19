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

func (s *Service) AcceptFriendRequest(ctx context.Context, receiver *models.AuthUser, senderId uint) (*models.AuthUser, error) {
	sender, err := s.datastore.GetUserById(ctx, senderId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	_, err = s.neo4jService.AcceptFriendRequest(ctx, sender, receiver)
	if err != nil {
		return nil, err
	}
	return receiver, nil
}
