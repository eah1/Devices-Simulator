// Package store contains methods to call database and initialise store.
package store

import (
	"device-simulator/business/core/models"
	"device-simulator/business/sys/db"
	"fmt"
)

// DeviceConfigCreate create a new device config.
func (s *Store) DeviceConfigCreate(deviceConfig models.DeviceConfig) error {
	if _, err := s.engine.Insert(deviceConfig); err != nil {
		return fmt.Errorf("store.device_config.DeviceConfigCreate.Insert(%+v): %w",
			deviceConfig, db.PsqlError(s.log, err))
	}

	if err := s.engine.Table(deviceConfig.TableName()).Commit(); err != nil {
		return fmt.Errorf("store.device_config.DeviceConfigCreate.Table(%+v): %w",
			deviceConfig, db.PsqlError(s.log, err))
	}

	return nil
}
