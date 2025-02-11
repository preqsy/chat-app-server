package models

import "time"

type AuthUser struct {
	Id        string    `json:"id" gorm:"primary_key"`
	Email     string    `json:"email" gorm:"unique:not null"`
	Password  string    `json:"password" gorm:"not null"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Username  string    `json:"username" gorm:"not null: unique"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type AuthUserCreate struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}
