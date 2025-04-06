package core

import (
	models "chat_app_server/model"
	"context"
)

func (s *Service) SaveMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	message, err := s.datastore.SaveMessage(ctx, message)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *Service) RetrieveMessages(ctx context.Context, senderId, receiverId int32) ([]*models.Message, error) {
	messages, err := s.datastore.RetrieveMessagesById(ctx, senderId, receiverId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *Service) RecentChats(ctx context.Context, senderId int32) ([]*models.Message, error) {
	recentChats, err := s.datastore.RecentChats(ctx, senderId)
	if err != nil {
		return nil, err
	}
	return recentChats, nil
}
