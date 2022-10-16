package main

import (
	"device-simulator/app/config"
	"device-simulator/app/services/myc-devices-simulator/handlers"
	"device-simulator/business/sys/handler"
	"device-simulator/business/sys/logger"
	"device-simulator/business/sys/sentry"
	"device-simulator/business/web/middlewares/common"
	"fmt"
	"os"
	"time"

	goSentry "github.com/getsentry/sentry-go"
	echoSentry "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	// Create config env vars.
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	// Construct the application logger.
	log, err := logger.InitLogger("MYC-DEVICES-SIMULATOR", cfg.Environment)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	fmt.Println(
		"starting environments status",
		fmt.Sprintf("host: %s", cfg.Host),
		fmt.Sprintf("hostName: %s", cfg.HostName),
		fmt.Sprintf("port: %s", cfg.Port),
		fmt.Sprintf("base url: %s", cfg.BaseURL),
		fmt.Sprintf("server URI: %s", cfg.ServerURI),
		fmt.Sprintf("environment: %s", cfg.Environment))

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

	// start services.
	log.Errorf("%s", startEcho(log, cfg))

	return nil
}

// startEcho start server.
func startEcho(log *zap.SugaredLogger, cfg config.Config) error {
	// Start App
	app := echo.New()

	// hide echo banner.
	app.HideBanner = true

	// Set logging level to INFO.
	app.Logger.SetLevel(2)

	// Config sentry echo.
	app.Use(echoSentry.New(echoSentry.Options{
		Repanic:         true,
		WaitForDelivery: false,
		Timeout:         time.Second,
	}))

	// aggregate common middlewares.
	common.AddCommonMiddlewares(app, log)

	// Initializing handles.
	handlerConfig := handler.NewHandlerConfig(cfg, log)
	handlers.Handlers(app, handlerConfig)

	return app.Start(cfg.Host + ":" + cfg.Port)
}
