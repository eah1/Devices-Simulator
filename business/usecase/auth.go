package usecase

import (
	"device-simulator/business/sys/auth"
	"device-simulator/business/web/webmodels"
)

// Login check credentials and user validated. Created token auth.
func (u *UseCase) Login(userLogin webmodels.LoginUser) (string, error) {
	user, err := u.core.User.FindByEmail(userLogin.Username)
	if err != nil {
		u.log.Errorw("Login error -> FindByEmail",
			"service", "USE CASE | Login | CORE USER", "error", err.Error())

		return "", err
	}

	if err := u.core.User.CheckCredentials(user, userLogin.Password); err != nil {
		u.log.Errorw("Login error -> CheckCredentials",
			"service", "USE CASE | Login | CORE USER", "error", err.Error())

		return "", err
	}

	if err := u.core.User.IsActivate(user); err != nil {
		u.log.Errorw("Login error -> IsActivate",
			"service", "USE CASE | Login | CORE USER", "error", err.Error())

		return "", err
	}

	clams := auth.CustomClaims{
		StandardClaims: auth.NewStandardClaims(),
		Email:          user.Email,
		ID:             user.ID,
	}

	token, err := u.core.Auth.GenerateToken(clams)
	if err != nil {
		u.log.Errorw("Login error -> GenerateToken",
			"service", "USE CASE | Login | CORE AUTH", "error", err.Error())

		return "", err
	}

	return token, nil
}
