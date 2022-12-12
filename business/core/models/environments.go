// Package models contains structs of the model application.
package models

import "time"

// Environment represents the structure we need for moving data.
type Environment struct {
	ID        string                 `xorm:"pk not null 'id'" json:"id"`
	UserID    string                 `xorm:"user_id" json:"userId"`
	Name      string                 `xorm:"name" json:"name"`
	Vars      map[string]interface{} `xorm:"jsonb vars" json:"vars"`
	CreatedAt time.Time              `xorm:"created" json:"-"`
	UpdatedAt time.Time              `xorm:"updated" json:"-"`
}
