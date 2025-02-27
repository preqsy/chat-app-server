package database

import (
	"fmt"

	models "chat_app_server/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostgresDB struct {
	client *gorm.DB
}

func (db *PostgresDB) SaveUser(user *models.AuthUser) (*models.AuthUser, error) {
	result := db.client.Create(user)
	if result.Error != nil {
		logrus.Error("Failed to save user: ", result.Error)
		return nil, result.Error
	}
	fmt.Println("User saved: ", user)
	return user, nil
}

func (db *PostgresDB) GetUserByEmail(email string) (*models.AuthUser, error) {
	var user models.AuthUser
	result := db.client.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil

}

func (db *PostgresDB) GetUserById(userId uint) (*models.AuthUser, error) {
	var user models.AuthUser
	result := db.client.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
