// Package usecase contains the logic of use cases.
package usecase

import (
	"device-simulator/business/core/models"
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
