// Package usecase contains the logic of use cases.
package usecase

import (
	"device-simulator/app/config"
	"device-simulator/business/core"
	"device-simulator/business/db/store"
	"device-simulator/business/web/webmodels"

	"go.uber.org/zap"
)

// UseCase build use case group.
type UseCase struct {
	log  *zap.SugaredLogger
	core core.Core
}

// NewUseCase constructs a use case group.
func NewUseCase(log *zap.SugaredLogger, config config.Config, store store.Store) UseCase {
	return UseCase{
		log:  log,
		core: core.NewCore(log, config, store),
	}
}

//go:generate mockery --name User --structname UserMock --filename UserMock.go

// User methods user use case.
type User interface {
	RegisterUser(userRegister webmodels.RegisterUser) error
}
