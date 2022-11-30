// Package core contains core business API.
package core

import (
	"device-simulator/app/config"
	"device-simulator/business/core/models"
	"device-simulator/business/db/store"
	errors2 "device-simulator/business/sys/errors"
	"device-simulator/foundation"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

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
		c.log.Errorw("GeneratePassword error",
			"service", "CORE | USER", "error", err.Error())

		return errors.Wrapf(err, "generating password hash")
	}

	user.Password = string(hash)

	return nil
}

// CheckCredentials compare password hash from password in database.
func (c *UserCore) CheckCredentials(user models.User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors2.ErrAuthenticationFailed
	}

	return nil
}

// Create insert a new user into the system.
func (c *UserCore) Create(user models.User) error {
	user.ID = uuid.New().String()

	if err := c.store.UserCreate(user); err != nil {
		c.log.Errorw("UserCreate error",
			"service", "CORE | USER", "error", err.Error())

		return err
	}

	return nil
}

// CreateValidationToken generate validation token from email activation.
func (c *UserCore) CreateValidationToken(user *models.User) error {
	validationToken, err := foundation.GenerateToken(16)
	if err != nil {
		c.log.Errorw("Create validation token",
			"service", "CORE | USER", "error", err.Error())

		return err
	}

	user.ValidationToken = *validationToken

	if err := c.store.UserUpdate(*user); err != nil {
		c.log.Errorw("CreateValidationToken error to Update user",
			"service", "CORE | USER", "error", err.Error())

		user.ValidationToken = ""

		return err
	}

	return nil
}

// Activate user activate in the system.
func (c *UserCore) Activate(user *models.User) error {
	if user.Validated {
		return errors2.ErrAuthenticationFailed
	}

	user.Validated = true

	if err := c.store.UserUpdate(*user); err != nil {
		c.log.Errorw("Activate error to Update user",
			"service", "CORE | USER", "error", err.Error())

		user.Validated = false

		return err
	}

	return nil
}

// IsActivate check user is activate.
func (c *UserCore) IsActivate(user models.User) error {
	if !user.Validated {
		return errors2.ErrAuthenticationFailed
	}

	return nil
}

// FindByEmail search user by email field.
func (c *UserCore) FindByEmail(email string) (models.User, error) {
	user, err := c.store.UserFindByEmail(email)
	if err != nil {
		c.log.Errorw("FindByEmail error",
			"service", "CORE | USER", "error", err.Error())

		return user, err
	}

	return user, nil
}

// FindByValidationToken search user by validation token field.
func (c *UserCore) FindByValidationToken(validationToken string) (models.User, error) {
	user, err := c.store.UserFindByValidationToken(validationToken)
	if err != nil {
		c.log.Errorw("FindByValidationToken error",
			"service", "CORE | USER", "error", err.Error())

		return user, err
	}

	return user, nil
}
