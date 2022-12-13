// Package store contains methods to call database and initialise store.
package store

import (
	"device-simulator/business/core/models"
	"device-simulator/business/sys/db"
	"fmt"
)

// EnvironmentCreate create a new environment.
func (s *Store) EnvironmentCreate(environment models.Environment) error {
	if _, err := s.engine.Insert(environment); err != nil {
		return fmt.Errorf("store.environment.EnvironmentCreate.Insert(%+v): %w",
			environment, db.PsqlError(s.log, err))
	}

	if err := s.engine.Table(environment.TableName()).Commit(); err != nil {
		return fmt.Errorf("store.environment.EnvironmentCreate.Table(%+v): %w",
			environment, db.PsqlError(s.log, err))
	}

	return nil
}
