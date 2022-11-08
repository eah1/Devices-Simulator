// Package task contains task queue business API.
package task

import (
	"encoding/json"
	"fmt"

	"device-simulator/business/task/models"
	"github.com/hibiken/asynq"
)

const (
	TypeSendValidationEmail = "sendValidationEmail"
)

// SendValidationEmail task send validation email and trigger async task.
func SendValidationEmail(emails, validationToken, language string) (*asynq.Task, error) {
	payload, err := json.Marshal(models.SendValidationEmail{
		Email: emails, ValidationToken: validationToken, Language: language,
	})
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return asynq.NewTask(TypeSendValidationEmail, payload, asynq.Queue("emails")), nil
}
