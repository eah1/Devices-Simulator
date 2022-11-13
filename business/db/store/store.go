// Package store contains methods to call database and initialise store.
package store

import (
	"device-simulator/business/core/models"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

// Store build store group.
type Store struct {
	log    *zap.SugaredLogger
	engine *xorm.Engine
}

// NewStore constructs a store group.
func NewStore(log *zap.SugaredLogger, engine *xorm.Engine) Store {
	return Store{
		log:    log,
		engine: engine,
	}
}

type SQL interface {
	UserCreate(user models.User) error
	UserFindByEmail(email string) (models.User, error)
	UserFindByValidationToken(validationToken string) (models.User, error)
	UserUpdate(user models.User) error
}
