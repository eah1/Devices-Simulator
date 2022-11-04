// Package store contains methods to call database and initialise store.
package store

import (
	"device-simulator/business/core/models"
	"device-simulator/business/sys/db"
)

// UserCreate creation a new user.
func (s *Store) UserCreate(user models.User) error {
	if _, err := s.engine.Insert(user); err != nil {
		return db.TranslatePsqlError(s.log, err)
	}

	if err := s.engine.Table(user.TableName()).Commit(); err != nil {
		return db.TranslatePsqlError(s.log, err)
	}

	return nil
}
