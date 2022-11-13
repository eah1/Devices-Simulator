// Package usecase contains the logic of use cases.
package usecase

import (
	"device-simulator/business/core/models"
	"device-simulator/business/task"
	"device-simulator/business/web/webmodels"
)

// RegisterUser register user case-use.
func (u *UseCase) RegisterUser(userRegister webmodels.RegisterUser) error {
	user := models.RegisterUserWebToUser(userRegister)

	if err := u.core.User.GeneratePassword(userRegister.Password, &user); err != nil {
		u.log.Errorw("RegisterUser error -> GeneratePassword",
			"service", "USE CASE | RegisterUser | CORE USER", "error", err.Error())

		return err
	}

	if err := u.core.User.Create(user); err != nil {
		u.log.Errorw("RegisterUser error -> Create",
			"service", "USE CASE | RegisterUser | CORE USER", "error", err.Error())

		return err
	}

	return nil
}

// SendValidationEmail sending email to validation account case-use.
func (u *UseCase) SendValidationEmail(email string) error {
	user, err := u.core.User.FindByEmail(email)
	if err != nil {
		u.log.Errorw("SendValidationEmail error -> FindByEmail",
			"service", "USE CASE | SendValidationEmail | CORE USER", "error", err.Error())
	}

	if err := u.core.User.CreateValidationToken(&user); err != nil {
		u.log.Errorw("SendValidationEmail error -> CreateValidationToken",
			"service", "USE CASE | SendValidationEmail | CORE USER", "error", err.Error())

		return err
	}

	sendValidation, err := task.SendValidationEmail(user.Email, user.ValidationToken, user.Language)
	if err != nil {
		u.log.Errorw("SendValidationEmail error -> SendValidationEmail",
			"service", "USE CASE | SendValidationEmail | TASK", "error", err.Error())
	}

	if _, err := u.clientQueue.Enqueue(sendValidation); err != nil {
		u.log.Errorw("SendValidationEmail error -> Enqueue task",
			"service", "USE CASE | SendValidationEmail | QUEUE", "error", err.Error())
	}

	return nil
}

// ActivateUser activation user case-use.
func (u *UseCase) ActivateUser(activateToken string) error {
	user, err := u.core.User.FindByValidationToken(activateToken)
	if err != nil {
		u.log.Errorw("ActivateUser error -> FindByValidationToken",
			"service", "USE CASE | ActivateUser | CORE USER", "error", err.Error())
	}

	if err := u.core.User.Activate(&user); err != nil {
		u.log.Errorw("ActivateUser error -> Activate",
			"service", "USE CASE | SendValidationEmail | CORE USER", "error", err.Error())

		return err
	}

	return nil
}
