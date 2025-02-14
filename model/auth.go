package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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

type AuthUserRegisterResponse struct {
	AuthUser AuthUser `json:"authUser"`
	Token    string   `json:"token"`
}

func (a *AuthUser) Validate() error {
	if err := validation.ValidateStruct(
		a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Username, validation.Required, validation.Length(3, 32)),
		validation.Field(&a.Password, validation.Required, validation.Length(4, 32)),
	); err != nil {
		return err
	}
	return nil
}
