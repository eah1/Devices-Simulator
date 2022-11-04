// Package handler contains the Handler configuration to app.
package handler

import (
	"device-simulator/app/config"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

// HandlerConfig contains all the mandatory systems required by handlers.
type HandlerConfig struct {
	Config config.Config
	Log    *zap.SugaredLogger
	DB     *xorm.Engine
}

// NewHandlerConfig initialize HandlerConfig structure.
func NewHandlerConfig(config config.Config, log *zap.SugaredLogger, db *xorm.Engine) HandlerConfig {
	handlerConfig := new(HandlerConfig)
	handlerConfig.Config = config
	handlerConfig.Log = log
	handlerConfig.DB = db

	return *handlerConfig
}
