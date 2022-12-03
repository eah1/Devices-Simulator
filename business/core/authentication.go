// Package core contains core business API.
package core

import (
	"device-simulator/app/config"
	"device-simulator/business/core/models"
	"device-simulator/business/db/store"
	mycErrors "device-simulator/business/sys/errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AuthenticationCore manages the set of API for authentication access.
type AuthenticationCore struct {
	log    *zap.SugaredLogger
	config config.Config
	store  store.Store
	core   *Core
}

// NewAuthenticationCore constructs a core for authentication API access.
func NewAuthenticationCore(
	log *zap.SugaredLogger, config config.Config, store store.Store, core *Core,
) AuthenticationCore {
	return AuthenticationCore{
		log:    log,
		config: config,
		store:  store,
		core:   core,
	}
}

// Create insert a new authentication into the system.
func (c *AuthenticationCore) Create(authentication models.Authentication) error {
	authentication.ID = uuid.New().String()

	if err := c.store.AuthenticationCreate(authentication); err != nil {
		return fmt.Errorf("core.authentication.create: %w", err)
	}

	return nil
}

// FindByTokenAndUserID search authentication by token and userId.
func (c *AuthenticationCore) FindByTokenAndUserID(token, userID string) (models.Authentication, error) {
	authentication, err := c.store.AuthenticationFindByTokenAndUserID(token, userID)
	if err != nil {
		return authentication, fmt.Errorf("core.authentication.FindByTokenAndUserID: %w", err)
	}

	return authentication, nil
}

// IsValid check authentication is valid.
func (c *AuthenticationCore) IsValid(authentication models.Authentication) error {
	if !authentication.Valid {
		return fmt.Errorf("core.authentication.IsValid: %w", mycErrors.ErrAuthenticationFailed)
	}

	return nil
}

// Invalidation activation flag valid to false.
func (c *AuthenticationCore) Invalidation(authentication *models.Authentication) error {
	if !authentication.Valid {
		return fmt.Errorf("core.authentication.Invalidation: %w", mycErrors.ErrAuthenticationFailed)
	}

	authentication.Valid = false
	authentication.LogoutAt = time.Now()

	if err := c.store.AuthenticationUpdate(*authentication); err != nil {
		authentication.Valid = true

		return fmt.Errorf("core.authentication.Invalidation: %w", err)
	}

	return nil
}
