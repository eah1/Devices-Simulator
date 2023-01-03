// Package webmodels contains core business API.
package webmodels

// CreateDevice  contains information needed to create a new Device.
type CreateDevice struct {
	Name           string `json:"name"           conform:"trim" validate:"required"`
	EnvironmentID  string `json:"environmentId"  conform:"trim" validate:"required,uuid"`
	DeviceConfigID string `json:"deviceConfigId" conform:"trim" validate:"required,uuid"`
}

// InformationDevice contains information.
type InformationDevice struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
