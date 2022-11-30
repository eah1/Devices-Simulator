// Package handlertasks contains handlers tasks from the queue.
package handlertasks

import (
	"context"
	"encoding/json"
	"fmt"

	"device-simulator/app/config"
	"device-simulator/business/core"
	"device-simulator/business/task/models"
	"github.com/hibiken/asynq"
	"github.com/jhillyerd/enmime"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

// MycSimulator build myc simulator struct.
type MycSimulator struct {
	Cfg         config.Config
	Log         *zap.SugaredLogger
	DB          *xorm.Engine
	Core        core.Core
	EmailSender *enmime.SMTPSender
}

// HandlerSendValidationEmail task process to send validation mail.
func (myc *MycSimulator) HandlerSendValidationEmail(_ context.Context, t *asynq.Task) error {
	var p models.SendValidationEmail

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	err := myc.Core.EmailSender.SendValidationEmail(p.Email, p.ValidationToken, p.Language)
	if err != nil {
		return err
	}

	return nil
}
