package test_test

import (
	"testing"

	"device-simulator/app/config"
	"device-simulator/business/sys/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// InitLogger create logger.
func InitLogger(t *testing.T, service string) *zap.SugaredLogger {
	t.Helper()

	log, err := logger.InitLogger(service, "local")
	require.NoError(t, err)

	return log
}

// InitConfig create a config env.
func InitConfig() config.Config {
	configENV := new(config.Config)

	return *configENV
}
