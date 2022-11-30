// Package models contains structs of the model application.
package models

import (
	"device-simulator/business/web/webmodels"
	"time"
)

// User represents the structure we need for moving data.
type User struct {
	ID              string    `xorm:"pk not null 'id'"      json:"id"`
	FirstName       string    `xorm:"first_name"            json:"firstName"`
	LastName        string    `xorm:"last_name"             json:"lastName"`
	Email           string    `xorm:"not null unique email" json:"email"`
	Password        string    `xorm:"password"              json:"-"`
	Language        string    `xorm:"language"              json:"language"`
	Company         string    `xorm:"company"               json:"company"`
	Validated       bool      `xorm:"validated"             json:"-"`
	ValidationToken string    `xorm:"validation_token"      json:"-"`
	CreatedAt       time.Time `xorm:"created"               json:"-"`
	UpdatedAt       time.Time `xorm:"updated"               json:"-"`
}

func (*User) TableName() string {
	return "users"
}

// RegisterUserWebToUser converter struct web model webmodels.RegisterUser to User model.
func RegisterUserWebToUser(userRegister webmodels.RegisterUser) User {
	user := new(User)
	user.FirstName = userRegister.FirstName
	user.LastName = userRegister.LastName
	user.Email = userRegister.Email
	user.Company = userRegister.Company
	user.Language = userRegister.Language

	return *user
}
