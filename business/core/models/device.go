// Package models contains structs of the model application.
package models

import "time"

// Device represents the structure we need for moving data.
type Device struct {
	ID             string    `xorm:"pk not null 'id'" json:"id"`
	Name           string    `xorm:"name"             json:"name"`
	UserID         string    `xorm:"user_id"          json:"userId"`
	EnvironmentID  string    `xorm:"environment_id"   json:"environmentId"`
	DeviceConfigID string    `xorm:"device_config_id" json:"deviceConfigId"`
	CreatedAt      time.Time `xorm:"created"          json:"-"`
	UpdatedAt      time.Time `xorm:"updated"          json:"-"`
}
