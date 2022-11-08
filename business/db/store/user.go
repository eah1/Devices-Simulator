// Package store contains methods to call database and initialise store.
package store

import (
	"device-simulator/business/core/models"
	"device-simulator/business/sys/db"
	errors2 "device-simulator/business/sys/errors"
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

// UserFindByEmail search user by email field.
func (s *Store) UserFindByEmail(email string) (models.User, error) {
	var user models.User

	res, err := s.engine.Where("email = ?", email).Get(&user)
	if err != nil {
		return user, db.TranslatePsqlError(s.log, err)
	}

	if !res {
		return user, errors2.ErrElementNotExist
	}

	return user, nil
}

// UserUpdate updates user.
func (s *Store) UserUpdate(user models.User) error {
	update, err := s.engine.UseBool().ID(user.ID).Update(&user)
	if err != nil {
		return db.TranslatePsqlError(s.log, err)
	}

	if update != 1 {
		return errors2.ErrElementNotExist
	}

	return nil
}
