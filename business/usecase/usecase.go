// Package usecase contains the logic of use cases.
package usecase

import (
	"device-simulator/app/config"
	"device-simulator/business/core"
	"device-simulator/business/db/store"
	"device-simulator/business/web/webmodels"
	"github.com/hibiken/asynq"
	"github.com/jhillyerd/enmime"
	"go.uber.org/zap"
)

// UseCase build use case group.
type UseCase struct {
	log         *zap.SugaredLogger
	core        core.Core
	clientQueue *asynq.Client
}

// NewUseCase constructs a use case group.
func NewUseCase(
	log *zap.SugaredLogger, config config.Config, store store.Store, clientQueue *asynq.Client,
	emailSender *enmime.SMTPSender,
) UseCase {
	return UseCase{
		log:         log,
		core:        core.NewCore(log, config, store, emailSender),
		clientQueue: clientQueue,
	}
}

//go:generate mockery --name User --structname UserMock --filename UserMock.go

// User methods user use case.
type User interface {
	RegisterUser(userRegister webmodels.RegisterUser) error
	SendValidationEmail(email string) error
	ActivateUser(activateToken string) error
}

// Auth methods auth use case.
type Auth interface {
	Login(userLogin webmodels.LoginUser) (string, error)
}
