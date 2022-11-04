// Package usecase contains the logic of use cases.
package usecase

import (
	"device-simulator/app/config"
	"device-simulator/business/core"
	"device-simulator/business/db/store"

	"go.uber.org/zap"
)

type UseCase struct {
	log  *zap.SugaredLogger
	core core.Core
}

func NewUseCase(log *zap.SugaredLogger, config config.Config, store store.Store) UseCase {
	return UseCase{
		log:  log,
		core: core.NewCore(log, config, store),
	}
}
