// Package store contains methods to call database and initialise store.
package store

import (
	"device-simulator/business/core/models"
	"device-simulator/business/sys/db"
	mycErrors "device-simulator/business/sys/errors"
	"fmt"
)

// AuthenticationCreate create a new authentication.
func (s *Store) AuthenticationCreate(authentication models.Authentication) error {
	if _, err := s.engine.Insert(authentication); err != nil {
		return fmt.Errorf("store.authentication.AuthenticationCreate.Insert(%+v): %w",
			authentication, db.PsqlError(s.log, err))
	}

	if err := s.engine.Table(authentication.TableName()).Commit(); err != nil {
		return fmt.Errorf("store.authentication.AuthenticationCreate.Table(%+v): %w",
			authentication, db.PsqlError(s.log, err))
	}

	return nil
}

// AuthenticationFindByTokenAndUserID search authentication by token and userId.
func (s *Store) AuthenticationFindByTokenAndUserID(token, userID string) (models.Authentication, error) {
	var authentication models.Authentication

	res, err := s.engine.Where("token = ? AND user_id = ?", token, userID).Get(&authentication)
	if err != nil {
		return authentication, fmt.Errorf("store.authentication.AuthenticationFindByTokenAndUserID(%s, %s): %w",
			token, userID, db.PsqlError(s.log, err))
	}

	if !res {
		return authentication, fmt.Errorf("store.authentication.AuthenticationFindByTokenAndUserID(%s, %s): %w",
			token, userID, mycErrors.ErrElementNotExist)
	}

	return authentication, nil
}

// AuthenticationUpdate updates an authentication.
func (s *Store) AuthenticationUpdate(authentication models.Authentication) error {
	res, err := s.engine.UseBool().ID(authentication.ID).Update(&authentication)
	if err != nil {
		return fmt.Errorf("store.authentication.AuthenticationUpdate(%+v): %w",
			authentication, db.PsqlError(s.log, err))
	}

	if res != 1 {
		return fmt.Errorf("store.authentication.AuthenticationUpdate(%+v): %w",
			authentication, mycErrors.ErrElementNotExist)
	}

	return nil
}
