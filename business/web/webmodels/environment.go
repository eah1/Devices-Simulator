// Package webmodels contains core business API.
package webmodels

// CreateEnvironment contains information needed to create a new Environment.
type CreateEnvironment struct {
	Name string             `json:"name" conform:"trim" validate:"required"`
	Vars []*EnvironmentVars `json:"vars"                validate:"required,min=1,dive"`
}

// EnvironmentVars fields to complete the CreateEnvironment struct.
type EnvironmentVars struct {
	Key string `json:"key" conform:"trim" validate:"required"`
	Var string `json:"var" conform:"trim" validate:"required"`
}

// InformationEnvironment contains information.
type InformationEnvironment struct {
	ID   string                 `json:"id"`
	Name string                 `json:"name"`
	Vars map[string]interface{} `json:"vars" swaggertype:"object,string" example:"key:value,key2:value2"`
}
