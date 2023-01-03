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

// DeviceCore manages the set of API for device core access.
type DeviceCore struct {
	log    *zap.SugaredLogger
	config config.Config
	store  store.Store
	core   *Core
}

// NewDeviceCore constructs a core for device API access.
func NewDeviceCore(log *zap.SugaredLogger, config config.Config, store store.Store, core *Core) DeviceCore {
	return DeviceCore{
		log:    log,
		config: config,
		store:  store,
		core:   core,
	}
}

// Create insert a new device into the system.
func (c *DeviceCore) Create(device *models.Device) error {
	device.ID = uuid.NewString()

	if err := c.store.DeviceCreate(*device); err != nil {
		return fmt.Errorf("core.device.create: %w", err)
	}

	return nil
}
