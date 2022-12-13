// Package models contains structs of the model application.
package models

import (
	"device-simulator/business/web/webmodels"
	"time"
)

// Environment represents the structure we need for moving data.
type Environment struct {
	ID        string                 `xorm:"pk not null 'id'" json:"id"`
	UserID    string                 `xorm:"user_id" json:"userId"`
	Name      string                 `xorm:"name" json:"name"`
	Vars      map[string]interface{} `xorm:"jsonb vars" json:"vars"`
	CreatedAt time.Time              `xorm:"created" json:"-"`
	UpdatedAt time.Time              `xorm:"updated" json:"-"`
}

func (*Environment) TableName() string {
	return "environments"
}

// CreateEnvironmentWebToEnvironment converter struct web model webmodels.CreateEnvironment to Environment model.
func CreateEnvironmentWebToEnvironment(createEnvironment webmodels.CreateEnvironment) Environment {
	environment := new(Environment)
	environment.Name = createEnvironment.Name

	environment.Vars = make(map[string]interface{})

	for _, vars := range createEnvironment.Vars {
		environment.Vars[vars.Key] = vars.Var
	}

	return *environment
}
