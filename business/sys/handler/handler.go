// Package handler contains the Handler configuration to app.
package handler

import (
	"device-simulator/app/config"

	"github.com/hibiken/asynq"
	"github.com/jhillyerd/enmime"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

// HandlerConfig contains all the mandatory systems required by handlers.
type HandlerConfig struct {
	Config            config.Config
	Log               *zap.SugaredLogger
	DB                *xorm.Engine
	ClientQueue       *asynq.Client
	EmailSender       *enmime.SMTPSender
	JWTConfig         echo.MiddlewareFunc
	AuthorizationUser echo.MiddlewareFunc
}

// NewHandlerConfig initialize HandlerConfig structure.
func NewHandlerConfig(
	config config.Config, log *zap.SugaredLogger, db *xorm.Engine, clientQueue *asynq.Client,
	emailSender *enmime.SMTPSender,
) HandlerConfig {
	handlerConfig := new(HandlerConfig)
	handlerConfig.Config = config
	handlerConfig.Log = log
	handlerConfig.DB = db
	handlerConfig.ClientQueue = clientQueue
	handlerConfig.EmailSender = emailSender

	return *handlerConfig
}
