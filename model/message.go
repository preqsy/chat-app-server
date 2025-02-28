package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model

	SenderID   uint   `json:"sender" gorm:"not null"`
	ReceiverID uint   `json:"receiver" gorm:"not null"`
	Content    string `json:"content" gorm:"not null"`

	Sender   AuthUser `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	Receiver AuthUser `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE"`
}
