// Package models contains structs of the model application.
package models

import (
	"time"
)

// Authentication represents the structure we need for moving data.
type Authentication struct {
	ID        string    `xorm:"pk not null 'id'" json:"id"`
	Token     string    `xorm:"token" json:"token"`
	UserID    string    `xorm:"user_id" json:"userId"`
	Valid     bool      `xorm:"valid" json:"valid"`
	LoginAt   time.Time `xorm:"login_at" json:"loginAt"`
	LogoutAt  time.Time `xorm:"logout_at" json:"logoutAt"`
	CreatedAt time.Time `xorm:"created" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
}

func (*Authentication) TableName() string {
	return "authentications"
}

// AuthenticationByToken create authentication model from token.
func AuthenticationByToken(token, userID string) Authentication {
	authentication := new(Authentication)
	authentication.Token = token
	authentication.UserID = userID
	authentication.Valid = true
	authentication.LoginAt = time.Now()

	return *authentication
}
