package database

import (
	models "chat_app_server/model"
	"context"

	"github.com/sirupsen/logrus"
)

func (db *PostgresDB) SaveMessage(ctx context.Context, message *models.Message) (*models.Message, error) {

	result := db.client.WithContext(ctx).Create(&message)

	if result.Error != nil {
		logrus.Error("Failed to save message: ", result.Error)
		return nil, result.Error
	}
	db.client.Preload("Sender").Preload("Receiver").First(&message, message.ID)
	return message, nil
}

func (db *PostgresDB) RetrieveMessagesById(ctx context.Context, senderId, receiverId int32) ([]*models.Message, error) {
	var messages []*models.Message
	err := db.client.WithContext(ctx).Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderId, receiverId, receiverId, senderId).Find(&messages).Error

	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (db *PostgresDB) RecentChats(ctx context.Context, userId int32) ([]*models.Message, error) {
	var recentChats []*models.Message

	subQuery := db.client.Table("messages").
		Select("MAX(id)").
		Where("sender_id = ? OR receiver_id = ?", userId, userId).
		Group("LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id)")

	err := db.client.WithContext(ctx).
		Preload("Sender").
		Preload("Receiver").
		Where("id IN (?)", subQuery).
		Order("id DESC").
		Find(&recentChats).Error

	if err != nil {
		return nil, err
	}

	return recentChats, nil
}
