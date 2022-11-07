// Package core contains core business API.
package core

import (
	"device-simulator/app/config"
	"device-simulator/business/db/store"

	"github.com/jhillyerd/enmime"
	"go.uber.org/zap"
)

// Core build core group.
type Core struct {
	User        UserCore
	EmailSender EmailSenderCore
}

// NewCore constructs a core group.
func NewCore(log *zap.SugaredLogger, config config.Config, store store.Store, emailSender *enmime.SMTPSender) Core {
	core := &Core{
		User:        NewUserCore(log, config, store, nil),
		EmailSender: NewEmailSenderCore(log, config, store, emailSender, nil),
	}

	return Core{
		User:        NewUserCore(log, config, store, core),
		EmailSender: NewEmailSenderCore(log, config, store, emailSender, core),
	}
}
