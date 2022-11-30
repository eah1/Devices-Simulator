// Package mycsimulator contains system myc library.
package mycsimulator

import (
	"device-simulator/app/config"
	"device-simulator/app/services/simulator-queue/handlertasks"
	"device-simulator/business/core"
	"device-simulator/business/db/store"

	"github.com/jhillyerd/enmime"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

// NewMycSimulator constructs a handlertasks.MycSimulator struct.
func NewMycSimulator(
	log *zap.SugaredLogger, cfg config.Config, database *xorm.Engine, emailSender *enmime.SMTPSender,
) handlertasks.MycSimulator {
	return handlertasks.MycSimulator{
		Log:         log,
		Cfg:         cfg,
		DB:          database,
		Core:        core.NewCore(log, cfg, store.NewStore(log, database), emailSender),
		EmailSender: emailSender,
	}
}
