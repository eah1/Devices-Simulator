package main

import (
	"device-simulator/app/config"
	"device-simulator/business/sys/db"
	"device-simulator/business/sys/emailsender"
	"device-simulator/business/sys/logger"
	"device-simulator/business/sys/mycsimulator"
	"device-simulator/business/sys/queue"
	"device-simulator/business/sys/sentry"
	"device-simulator/business/task"
	"fmt"
	"os"
	"strings"

	goSentry "github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

const retryDBMax = 5

func main() {
	// Create config env vars.
	cfg, err := config.LoadConfig()
	if err != nil {
		os.Exit(1)
	}

	// Construct the application logger.
	log, err := logger.InitLogger("SIMULATOR-QUEUE", cfg.Environment)
	if err != nil {
		os.Exit(1)
	}

	log.Infow("starting environments status",
		"host", cfg.Host, "hostName", cfg.HostName, "port", cfg.Port,
		"base url", cfg.BaseURL, "server URI", cfg.ServerURI, "environment", cfg.Environment)

	// Perform the startup and shutdown sequence.
	if err := run(log, cfg); err != nil {
		log.Errorw("startup", "ERROR", err)

		os.Exit(1)
	}

	defer func(log *zap.SugaredLogger) {
		if err := log.Sync(); err != nil {
			log.Error(err)
		}
	}(log)
}

// run init app.
func run(log *zap.SugaredLogger, cfg config.Config) error {
	log.Infow("startup")

	// Created a basic configuration sentry.
	if err := goSentry.Init(sentry.InitSentryConfig(cfg)); err != nil {
		log.Errorf("sentry configuration: %s", err)

		return fmt.Errorf("sentry configuration: %w", err)
	}

	log.Infow("starting sentry config status")

	// Create connectivity to the database.
	host := cfg.DBPostgres[strings.Index(cfg.DBPostgres, "@")+1 : strings.LastIndex(cfg.DBPostgres, "/")]

	database, err := db.Open(db.NewConfigDB(cfg), retryDBMax, log)
	if err != nil {
		log.Errorf("database open: %s", err)

		return fmt.Errorf("database open: %w", err)
	}

	log.Infow("starting database status", "host", host)

	// Created a configuration email.
	emailSender, err := emailsender.InnitEmailConfig(cfg)
	if err != nil {
		log.Errorf("email configuration: %w", err)

		return fmt.Errorf("email config: %w", err)
	}

	log.Infow("starting email sender status")

	// Create structure myc.
	myc := mycsimulator.NewMycSimulator(log, cfg, database, emailSender)

	app := asynq.NewServeMux()

	// Define task of the queue.
	app.HandleFunc(task.TypeSendValidationEmail, myc.HandlerSendValidationEmail)

	// Run process queue.
	if err := queue.NewQueue(cfg, log).Run(app); err != nil {
		log.Error("could not run server: %v", err)
	}

	return nil
}
