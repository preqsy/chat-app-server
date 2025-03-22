package core

import (
	models "chat_app_server/model"
	"context"
	"errors"
)

func (s *Service) SendFriendRequest(ctx context.Context, sender *models.AuthUser, receiverId uint) (*models.AuthUser, error) {
	// Fetch the receiver user
	receiver, err := s.datastore.GetUserById(ctx, receiverId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Prevent self-friend request
	if receiverId == sender.ID {
		return nil, errors.New("can't send a friend request to yourself")
	}

	isFriends, err := s.neo4jService.CheckIfFriends(ctx, sender, receiver)
	if err != nil {
		return nil, err
	}
	if isFriends {
		return nil, errors.New("you are already friends")
	}

	hasPendingRequest, err := s.neo4jService.CheckFriendRequest(ctx, sender, receiver)
	if err != nil {
		return nil, err
	}
	if hasPendingRequest {
		return nil, errors.New("friend request already sent or pending")
	}

	_, err = s.neo4jService.SendFriendRequest(ctx, sender, receiver)
	if err != nil {
		return nil, errors.New("failed to send friend request")
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

func (s *Service) ListFriendRequests(ctx context.Context, skip, limit int32, user *models.AuthUser) ([]*models.AuthUser, error) {
	result, err := s.neo4jService.ListFriendRequests(ctx, user)
	if err != nil {
		return nil, err
	}
	users, err := s.datastore.ListUsersByIds(ctx, result, skip, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (s *Service) ListFriends(ctx context.Context, skip, limit int32, user *models.AuthUser) ([]*models.AuthUser, error) {
	result, err := s.neo4jService.ListFriends(ctx, user)
	if err != nil {
		return nil, err
	}
	users, err := s.datastore.ListUsersByIds(ctx, result, skip, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}
