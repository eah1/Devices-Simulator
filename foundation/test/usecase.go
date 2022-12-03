package test_test

import (
	"device-simulator/business/db/store"
	"device-simulator/business/usecase"
	"device-simulator/business/web/webmodels"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func UseCaseRegisterValidate(
	t *testing.T, useCase usecase.UseCase, store store.Store, testName string,
) (string, string) {
	t.Helper()

	user := NewRegistrationUser(testName)
	assert.Nil(t, useCase.RegisterUser(user))

	// send email validation user.
	assert.Nil(t, useCase.SendValidationEmail(user.Email))

	// find user in database.
	userDB, err := store.UserFindByEmail(user.Email)
	require.NoError(t, err)

	assert.Nil(t, useCase.ActivateUser(userDB.ValidationToken))

	// find user in database.
	userDB, err = store.UserFindByEmail(user.Email)
	require.NoError(t, err)

	assert.True(t, userDB.Validated)

	return user.Email, user.Password
}

func UseCaseLogin(t *testing.T, useCase usecase.UseCase, email, password string) string {
	t.Helper()

	loginWebModel := webmodels.LoginUser{
		Username: email,
		Password: password,
	}

	token, err := useCase.Login(loginWebModel)

	assert.NotEmpty(t, token)
	assert.Nil(t, err)

	return token
}
