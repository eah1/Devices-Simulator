// Package models contains structs of the model application.
package models

import "time"

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
