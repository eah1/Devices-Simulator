package usecase

import (
	"device-simulator/business/sys/auth"
	"device-simulator/business/web/webmodels"
	"fmt"
)

// Login check credentials and user validated. Created token auth.
func (u *UseCase) Login(userLogin webmodels.LoginUser) (string, error) {
	user, err := u.core.User.FindByEmail(userLogin.Username)
	if err != nil {
		return "", fmt.Errorf("usecase.auth.Login.FindByEmail(%s): %w", userLogin.Username, err)
	}

	if err := u.core.User.CheckCredentials(user, userLogin.Password); err != nil {
		return "", fmt.Errorf("usecase.auth.Login.CheckCredentials(%+v, %s): %w",
			user, userLogin.Password, err)
	}

	if err := u.core.User.IsActivate(user); err != nil {
		return "", fmt.Errorf("usecase.auth.Login.IsActivate(%+v): %w", user, err)
	}

	clams := auth.CustomClaims{
		StandardClaims: auth.NewStandardClaims(),
		Email:          user.Email,
		ID:             user.ID,
	}

	token, err := u.core.Auth.GenerateToken(clams)
	if err != nil {
		return "", fmt.Errorf("usecase.auth.Login.GenerateToken(%+v): %w", clams, err)
	}

	return token, nil
}
