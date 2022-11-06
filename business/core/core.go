// Package core contains core business API.
package core

import (
	"device-simulator/app/config"
	"device-simulator/business/db/store"

	"go.uber.org/zap"
)

// Core build core group.
type Core struct {
	User UserCore
}

// NewCore constructs a core group.
func NewCore(log *zap.SugaredLogger, config config.Config, store store.Store) Core {
	core := &Core{
		User: NewUserCore(log, config, store, nil),
	}

	return Core{
		User: NewUserCore(log, config, store, core),
	}
}
