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
