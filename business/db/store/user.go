// Package store contains methods to call database and initialise store.
package store

import (
	"device-simulator/business/core/models"
	"device-simulator/business/sys/db"
	mycErrors "device-simulator/business/sys/errors"
	"fmt"
)

// UserCreate creation a new user.
func (s *Store) UserCreate(user models.User) error {
	if _, err := s.engine.Insert(user); err != nil {
		return fmt.Errorf("store.user.UserCreate.Insert(%+v): %w", user, db.PsqlError(s.log, err))
	}

	if err := s.engine.Table(user.TableName()).Commit(); err != nil {
		return fmt.Errorf("store.user.UserCreate.Table(%+v): %w", user, db.PsqlError(s.log, err))
	}

	return nil
}

// UserFindByID search user by id field.
func (s *Store) UserFindByID(userID string) (models.User, error) {
	user := new(models.User)

	res, err := s.engine.ID(userID).Get(user)
	if err != nil {
		return models.User{}, fmt.Errorf("store.user.UserFindByID(%s): %w", userID, db.PsqlError(s.log, err))
	}

	if !res {
		return models.User{}, fmt.Errorf("store.user.UserFindByID(%s): %w", userID, mycErrors.ErrElementNotExist)
	}

	return *user, nil
}

// UserFindByEmail search user by email field.
func (s *Store) UserFindByEmail(email string) (models.User, error) {
	var user models.User

	res, err := s.engine.Where("email = ?", email).Get(&user)
	if err != nil {
		return user, fmt.Errorf("store.user.UserFindByEmail(%s): %w", email, db.PsqlError(s.log, err))
	}

	if !res {
		return user, fmt.Errorf("store.user.UserFindByEmail(%s): %w", email, mycErrors.ErrElementNotExist)
	}

	return user, nil
}

// UserFindByValidationToken search user by validation token.
func (s *Store) UserFindByValidationToken(validationToken string) (models.User, error) {
	var user models.User

	res, err := s.engine.Where("validation_token = ?", validationToken).Get(&user)
	if err != nil {
		return user, fmt.Errorf("store.user.UserFindByValidationToken(%s): %w", validationToken, db.PsqlError(s.log, err))
	}

	if !res {
		return user, fmt.Errorf("store.user.UserFindByValidationToken(%s): %w", validationToken, mycErrors.ErrElementNotExist)
	}

	return user, nil
}

// UserUpdate updates user.
func (s *Store) UserUpdate(user models.User) error {
	update, err := s.engine.UseBool().ID(user.ID).Update(&user)
	if err != nil {
		return fmt.Errorf("store.user.UserUpdate.Update(%+v): %w", user, db.PsqlError(s.log, err))
	}

	if update != 1 {
		return fmt.Errorf("store.user.UserUpdate(%+v): %w", user, mycErrors.ErrElementNotExist)
	}

	return nil
}
