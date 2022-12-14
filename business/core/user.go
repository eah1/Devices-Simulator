// Package core contains core business API.
package core

import (
	"device-simulator/app/config"
	"device-simulator/business/core/models"
	"device-simulator/business/db/store"
	mycErrors "device-simulator/business/sys/errors"
	"device-simulator/foundation"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const sizeToken = 16

// UserCore manages the set of API for user access.
type UserCore struct {
	log    *zap.SugaredLogger
	config config.Config
	store  store.Store
	core   *Core
}

// NewUserCore constructs a core for user API access.
func NewUserCore(log *zap.SugaredLogger, config config.Config, store store.Store, core *Core) UserCore {
	return UserCore{
		log:    log,
		config: config,
		store:  store,
		core:   core,
	}
}

// GeneratePassword generate password hash.
func (c *UserCore) GeneratePassword(password string, user *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("core.user.GeneratePassword: %w", err)
	}

	user.Password = string(hash)

	return nil
}

// CheckCredentials compare password hash from password in database.
func (c *UserCore) CheckCredentials(user models.User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return fmt.Errorf("core.user.CheckCredentials: %w", mycErrors.ErrAuthenticationFailed)
	}

	return nil
}

// Create insert a new user into the system.
func (c *UserCore) Create(user models.User) error {
	user.ID = uuid.NewString()

	if err := c.store.UserCreate(user); err != nil {
		return fmt.Errorf("core.user.Create: %w", err)
	}

	return nil
}

// Update updates a user from system.
func (c *UserCore) Update(user models.User) error {
	if err := c.store.UserUpdate(user); err != nil {
		return fmt.Errorf("core.user.Update: %w", err)
	}

	return nil
}

// CreateValidationToken generate validation token from email activation.
func (c *UserCore) CreateValidationToken(user *models.User) error {
	validationToken, err := foundation.GenerateToken(sizeToken)
	if err != nil {
		return fmt.Errorf("core.user.CreateValidationToken: %w", err)
	}

	user.ValidationToken = *validationToken

	if err := c.store.UserUpdate(*user); err != nil {
		user.ValidationToken = ""

		return fmt.Errorf("core.user.CreateValidationToken: %w", err)
	}

	return nil
}

// Activate user activate in the system.
func (c *UserCore) Activate(user *models.User) error {
	if user.Validated {
		return fmt.Errorf("core.user.Activate: %w", mycErrors.ErrAuthenticationFailed)
	}

	user.Validated = true

	if err := c.store.UserUpdate(*user); err != nil {
		user.Validated = false

		return fmt.Errorf("core.user.Activate: %w", mycErrors.ErrAuthenticationFailed)
	}

	return nil
}

// IsActivate check user is activated.
func (c *UserCore) IsActivate(user models.User) error {
	if !user.Validated {
		return fmt.Errorf("core.user.IsActivate: %w", mycErrors.ErrAuthenticationFailed)
	}

	return nil
}

// FindByID search user by id field.
func (c *UserCore) FindByID(userID string) (models.User, error) {
	user, err := c.store.UserFindByID(userID)
	if err != nil {
		return user, fmt.Errorf("core.user.FindByID: %w", err)
	}

	return user, nil
}

// FindByEmail search user by email field.
func (c *UserCore) FindByEmail(email string) (models.User, error) {
	user, err := c.store.UserFindByEmail(email)
	if err != nil {
		return user, fmt.Errorf("core.user.FindByEmail: %w", err)
	}

	return user, nil
}

// FindByValidationToken search user by validation token field.
func (c *UserCore) FindByValidationToken(validationToken string) (models.User, error) {
	user, err := c.store.UserFindByValidationToken(validationToken)
	if err != nil {
		return user, fmt.Errorf("core.user.FindByValidationToken: %w", err)
	}

	return user, nil
}
