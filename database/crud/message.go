package database

import (
	models "chat_app_server/model"
	"context"

	"github.com/sirupsen/logrus"
)

func (db *PostgresDB) SaveMessage(ctx context.Context, message *models.Message) (*models.Message, error) {

	result := db.client.Create(&message)

	if result.Error != nil {
		logrus.Error("Failed to save user: ", result.Error)
		return nil, result.Error
	}
	return message, nil
}
