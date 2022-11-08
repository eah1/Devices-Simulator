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
	handlerConfig.DB = InitDatabase(t, handlerConfig.Config, handlerConfig.Log)

	handlerConfig.Config.TemplateFolder = "../../business/template/"

	handlerConfig.EmailSender = InitEmailConfig(t, handlerConfig.Config)
	handlerConfig.ClientQueue = InitClientQueue(t, handlerConfig.Config)

	return *handlerConfig
}
