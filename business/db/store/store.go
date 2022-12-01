// Package store contains methods to call database and initialise store.
package store

import (
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
