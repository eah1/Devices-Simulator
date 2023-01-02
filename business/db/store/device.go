// Package store contains methods to call database and initialise store.
package store

import (
	"device-simulator/business/core/models"
	"device-simulator/business/sys/db"
	"fmt"
)

// DeviceCreate create a new device.
func (s *Store) DeviceCreate(device models.Device) error {
	if _, err := s.engine.Insert(device); err != nil {
		return fmt.Errorf("store.device.DeviceCreate.Insert(%+v): %w",
			device, db.PsqlError(s.log, err))
	}

	if err := s.engine.Table(device.TableName()).Commit(); err != nil {
		return fmt.Errorf("store.device.DeviceCreate.Table(%+v): %w",
			device, db.PsqlError(s.log, err))
	}

	return nil
}
