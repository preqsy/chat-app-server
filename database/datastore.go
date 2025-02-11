package datastore

import (
	models "chat_app_server/model"
)

type Datastore interface {
	SaveUser(user *models.AuthUser) (*models.AuthUser, error)
}
