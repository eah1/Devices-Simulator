// Package handler contains the Handler configuration to app.
package handler

import (
	"device-simulator/app/config"

	"go.uber.org/zap"
)

// HandlerConfig contains all the mandatory systems required by handlers.
type HandlerConfig struct {
	Config config.Config
	Log    *zap.SugaredLogger
}

// NewHandlerConfig initialize HandlerConfig structure.
func NewHandlerConfig(config config.Config, log *zap.SugaredLogger) HandlerConfig {
	handlerConfig := new(HandlerConfig)
	handlerConfig.Config = config
	handlerConfig.Log = log

	return *handlerConfig
}
