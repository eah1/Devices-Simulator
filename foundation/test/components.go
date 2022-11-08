package test_test

import (
	"device-simulator/app/config"
	"device-simulator/business/sys/db"
	"device-simulator/business/sys/emailsender"
	"device-simulator/business/sys/logger"
	"os"
	"sync"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/jhillyerd/enmime"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

var (
	lock     = &sync.Mutex{}
	database *xorm.Engine
	queue    *asynq.Client
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
	configENV.QueueHost = "localhost"
	configENV.QueuePort = "6379"
	configENV.QueueConcurrency = 15
	configENV.PostmarkToken = "fakeMailerToken"
	configENV.SMTPFrom = "no-reply@circutor.com"
	configENV.SMTPHost = "localhost"
	configENV.SMTPPort = "25"
	configENV.SMTPNetwork = "tcp"

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

// InitClientQueue create queue client configuration.
func InitClientQueue(t *testing.T, config config.Config) *asynq.Client {
	t.Helper()

	if queue == nil {
		lock.Lock()
		defer lock.Unlock()

		if queue == nil {
			queue = asynq.NewClient(asynq.RedisClientOpt{Addr: config.QueueHost + ":" + config.QueuePort})
		}
	}

	return queue
}

// InitEmailConfig create email configuration.
func InitEmailConfig(t *testing.T, config config.Config) *enmime.SMTPSender {
	t.Helper()

	emailSender, err := emailsender.InnitEmailConfig(config)
	require.NoError(t, err)

	return emailSender
}
