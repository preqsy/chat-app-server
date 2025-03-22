package database

import (
	"context"
	"fmt"

	models "chat_app_server/model"

	"github.com/sirupsen/logrus"
)

func (db *PostgresDB) SaveUser(ctx context.Context, user *models.AuthUser) (*models.AuthUser, error) {
	result := db.client.Create(user)
	if result.Error != nil {
		logrus.Error("Failed to save user: ", result.Error)
		return nil, result.Error
	}
	fmt.Println("User saved: ", user)
	return user, nil
}

func (db *PostgresDB) GetUserByEmail(ctx context.Context, email string) (*models.AuthUser, error) {
	var user models.AuthUser
	result := db.client.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil

}

func (db *PostgresDB) GetUserById(ctx context.Context, userId uint) (*models.AuthUser, error) {
	var user models.AuthUser
	result := db.client.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (db *PostgresDB) ListUsers(ctx context.Context, skip, limit int32, id uint) ([]*models.AuthUser, error) {
	var users []*models.AuthUser
	result := db.client.Where("id <> ? ", id).Offset(int(skip)).Limit(int(limit)).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (db *PostgresDB) ListUsersByIds(ctx context.Context, ids []int64, skip, limit int32) ([]*models.AuthUser, error) {
	var users []*models.AuthUser
	result := db.client.Where("id IN (?)", ids).Offset(int(skip)).Limit(int(limit)).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
