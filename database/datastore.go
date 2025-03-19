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
	ListUsers(ctx context.Context, skip, limit int32) ([]*models.AuthUser, error)
}
