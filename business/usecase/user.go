// Package usecase contains the logic of use cases.
package usecase

import (
	"device-simulator/business/core/models"
	"device-simulator/business/task"
	"device-simulator/business/web/webmodels"
	"fmt"
)

// RegisterUser register user case-use.
func (u *UseCase) RegisterUser(userRegister webmodels.RegisterUser) error {
	user := models.RegisterUserWebToUser(userRegister)

	if err := u.core.User.GeneratePassword(userRegister.Password, &user); err != nil {
		return fmt.Errorf("usecase.user.RegisterUser.GeneratePassword(%s, %+v): %w", userRegister.Password, &user, err)
	}

	if err := u.core.User.Create(user); err != nil {
		return fmt.Errorf("usecase.user.RegisterUser.Create(%+v): %w", &user, err)
	}

	return nil
}

// SendValidationEmail sending email to validation account case-use.
func (u *UseCase) SendValidationEmail(email string) error {
	user, err := u.core.User.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("usecase.user.SendValidationEmail.FindByEmail(%s): %w", email, err)
	}

	if err := u.core.User.CreateValidationToken(&user); err != nil {
		return fmt.Errorf("usecase.user.SendValidationEmail.CreateValidationToken(%+v): %w", &user, err)
	}

	sendValidation, err := task.SendValidationEmail(user.Email, user.ValidationToken, user.Language)
	if err != nil {
		return fmt.Errorf("usecase.user.SendValidationEmail.task.SendValidationEmail(%s, %s, %s): %w",
			user.Email, user.ValidationToken, user.Language, err)
	}

	if _, err := u.clientQueue.Enqueue(sendValidation); err != nil {
		return fmt.Errorf("usecase.user.SendValidationEmail.Enqueue: %w", err)
	}

	return nil
}

// ActivateUser activation user case-use.
func (u *UseCase) ActivateUser(activateToken string) error {
	user, err := u.core.User.FindByValidationToken(activateToken)
	if err != nil {
		return fmt.Errorf("usecase.user.ActivateUser.FindByValidationToken(%s): %w", activateToken, err)
	}

	if err := u.core.User.Activate(&user); err != nil {
		return fmt.Errorf("usecase.user.ActivateUser.Activate(%+v): %w", &user, err)
	}

	return nil
}
