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

// DeviceConfigCore manages the set of API for device core access.
type DeviceConfigCore struct {
	log    *zap.SugaredLogger
	config config.Config
	store  store.Store
	core   *Core
}

// NewDeviceConfigCore constructs a core for device config API access.
func NewDeviceConfigCore(log *zap.SugaredLogger, config config.Config, store store.Store, core *Core) DeviceConfigCore {
	return DeviceConfigCore{
		log:    log,
		config: config,
		store:  store,
		core:   core,
	}
}

// Create insert a new device config into the system.
func (c *DeviceConfigCore) Create(deviceConfig *models.DeviceConfig) error {
	deviceConfig.ID = uuid.NewString()

	if err := c.store.DeviceConfigCreate(*deviceConfig); err != nil {
		return fmt.Errorf("core.device_config.create: %w", err)
	}

	return nil
}
