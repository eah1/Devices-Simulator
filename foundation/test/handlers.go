package test_test

import (
	"device-simulator/business/sys/handler"
	"testing"
)

// InitHandlerConfig create a config handler.
func InitHandlerConfig(t *testing.T, service string) handler.HandlerConfig {
	t.Helper()

	handlerConfig := new(handler.HandlerConfig)
	handlerConfig.Log = InitLogger(t, service)
	handlerConfig.Config = InitConfig()

	return *handlerConfig
}
