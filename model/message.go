package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model

	Sender   uint   `json:"sender" gorm:"not null"`
	Receiver uint   `json:"receiver" gorm:"not null"`
	Content  string `json:"content" gorm:"not null"`
}
