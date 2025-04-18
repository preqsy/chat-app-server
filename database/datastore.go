package datastore

import (
	models "chat_app_server/model"
	"context"
)

type Datastore interface {
	SaveUser(ctx context.Context, user *models.AuthUser) (*models.AuthUser, error)
	GetUserByEmail(ctx context.Context, email string) (*models.AuthUser, error)
	GetUserById(ctx context.Context, userId uint) (*models.AuthUser, error)
	SaveMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	ListUsers(ctx context.Context, skip, limit int32, ids []int64) ([]*models.AuthUser, error)
	ListUsersByIds(ctx context.Context, ids []int64, skip, limit int32) ([]*models.AuthUser, error)
	RetrieveMessagesById(ctx context.Context, senderId, receiverId int32) ([]*models.Message, error)
	RecentChats(ctx context.Context, senderId int32) ([]*models.Message, error)
}
