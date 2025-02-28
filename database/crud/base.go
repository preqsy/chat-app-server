package database

import (
	models "chat_app_server/model"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	client *gorm.DB
}

func ConnectDB(host, user, password, dbname, port string) (*PostgresDB, error) {
	logrus.Info("Connecting to database.........")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error("Failed to connect to database", err.Error())
		return nil, err
	}
	logrus.Info("Connected to database")
	db.AutoMigrate(&models.AuthUser{}, &models.Message{})

	return &PostgresDB{client: db}, nil
}
