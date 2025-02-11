package models

import (
	"gorm.io/gorm"
)

type AuthUser struct {
	gorm.Model

	// Id        uint   `json:"id" gorm:"primaryKey"`
	Email     string `json:"email" gorm:"unique:not null"`
	Password  string `json:"password" gorm:"not null"`
	FirstName string `json:"firstName" gorm:"not null"`
	LastName  string `json:"lastName" gorm:"not null"`
	Username  string `json:"username" gorm:"not null: unique"`
}

type AuthUserCreate struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
}
