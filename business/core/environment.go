// Package core contains core business API.
package core

import (
	"device-simulator/app/config"
	"device-simulator/business/core/models"
	"device-simulator/business/db/store"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// EnvironmentCore manages the set of API for environment access.
type EnvironmentCore struct {
	log    *zap.SugaredLogger
	config config.Config
	store  store.Store
	core   *Core
}

// NewEnvironmentCore constructs a core for environment API access.
func NewEnvironmentCore(log *zap.SugaredLogger, config config.Config, store store.Store, core *Core) EnvironmentCore {
	return EnvironmentCore{
		log:    log,
		config: config,
		store:  store,
		core:   core,
	}
}

// Create insert a new environment into the system.
func (c *EnvironmentCore) Create(environment *models.Environment) error {
	environment.ID = uuid.NewString()

	if err := c.store.EnvironmentCreate(*environment); err != nil {
		return fmt.Errorf("core.environment.create: %w", err)
	}

	return nil
}
