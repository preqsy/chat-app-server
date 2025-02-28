package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model

	Sender   string `json:"sender" gorm:"not null"`
	Receiver string `json:"receiver" gorm:"not null"`
	Content  string `json:"content" gorm:"not null"`
}
