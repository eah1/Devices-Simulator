// Package core contains core business API.
package core

import (
	"device-simulator/app/config"
	"device-simulator/business/core/models"
	"device-simulator/business/db/store"

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
