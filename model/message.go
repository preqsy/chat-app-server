package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model

	SenderID   uint   `json:"sender_id" gorm:"not null"`
	ReceiverID uint   `json:"receiver_id" gorm:"not null"`
	Content    string `json:"content" gorm:"not null"`

	Sender   AuthUser `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"-"`
	Receiver AuthUser `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE" json:"-"`
}
