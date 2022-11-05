package test_test

import (
	"device-simulator/app/config"
	"device-simulator/business/sys/db"
	"device-simulator/business/sys/logger"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

var (
	lock     = &sync.Mutex{}
	database *xorm.Engine
)

// InitLogger create logger.
func InitLogger(t *testing.T, service string) *zap.SugaredLogger {
	t.Helper()

	log, err := logger.InitLogger(service, "test")
	require.NoError(t, err)

	return log
}

// InitConfig create a config env.
func InitConfig() config.Config {
	configENV := new(config.Config)
	configENV.DBPostgres = os.Getenv("MYC_DEVICES_SIMULATOR_DBPOSTGRES")
	configENV.DBMaxOpenConns = 25
	configENV.DBMaxIdleConns = 25
	configENV.DBLogger = false

	return *configENV
}

// InitDatabase create database.
func InitDatabase(t *testing.T, config config.Config, log *zap.SugaredLogger) *xorm.Engine {
	t.Helper()

	var err error

	if database == nil {
		lock.Lock()
		defer lock.Unlock()

		if database == nil {
			database, err = db.Open(db.NewConfigDB(config), 5, log)
			require.NoError(t, err)
		}
	}

	return database
}
