package usecase_test

import (
	"device-simulator/business/db/store"
	errors2 "device-simulator/business/sys/errors"
	"device-simulator/business/usecase"
	"device-simulator/business/web/webmodels"
	tt "device-simulator/foundation/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestUseCaseLogin(t *testing.T) {
	t.Parallel()

	testName := "use-case-login"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	newUseCase := usecase.NewUseCase(
		newLog, newConfig, newStore, tt.InitClientQueue(t, newConfig), tt.InitEmailConfig(t, newConfig))

	t.Log("Given the need to work with the login user use case.")
	{
		t.Logf("\tWhen a correct login authorization.")
		{
			user := tt.NewRegistrationUser(testName)
			assert.Nil(t, newUseCase.RegisterUser(user))

			// send email validation user.
			assert.Nil(t, newUseCase.SendValidationEmail(user.Email))

			// find user in database.
			userDB, err := newStore.UserFindByEmail(user.Email)
			require.NoError(t, err)

			assert.Nil(t, newUseCase.ActivateUser(userDB.ValidationToken))

			// find user in database.
			userDB, err = newStore.UserFindByEmail(user.Email)
			require.NoError(t, err)

			assert.True(t, userDB.Validated)

			loginWebModel := webmodels.LoginUser{
				Username: user.Email,
				Password: user.Password,
			}

			token, err := newUseCase.Login(loginWebModel)

			assert.NotEmpty(t, token)
			assert.Nil(t, err)
		}

		t.Logf("\tWhen a failed login authoritzation when username not exist.")
		{
			user := tt.NewRegistrationUser(testName)
			assert.Nil(t, newUseCase.RegisterUser(user))

			// send email validation user.
			assert.Nil(t, newUseCase.SendValidationEmail(user.Email))

			// find user in database.
			userDB, err := newStore.UserFindByEmail(user.Email)
			require.NoError(t, err)

			assert.Nil(t, newUseCase.ActivateUser(userDB.ValidationToken))

			// find user in database.
			userDB, err = newStore.UserFindByEmail(user.Email)
			require.NoError(t, err)

			assert.True(t, userDB.Validated)

			loginWebModel := webmodels.LoginUser{
				Username: faker.Internet().Email(),
				Password: user.Password,
			}

			token, err := newUseCase.Login(loginWebModel)

			assert.Equal(t, "", token)
			assert.Error(t, errors2.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a failed login authoritzation when password not correct.")
		{
			user := tt.NewRegistrationUser(testName)
			assert.Nil(t, newUseCase.RegisterUser(user))

			// send email validation user.
			assert.Nil(t, newUseCase.SendValidationEmail(user.Email))

			// find user in database.
			userDB, err := newStore.UserFindByEmail(user.Email)
			require.NoError(t, err)

			assert.Nil(t, newUseCase.ActivateUser(userDB.ValidationToken))

			// find user in database.
			userDB, err = newStore.UserFindByEmail(user.Email)
			require.NoError(t, err)

			assert.True(t, userDB.Validated)

			loginWebModel := webmodels.LoginUser{
				Username: user.Email,
				Password: faker.RandomString(20),
			}

			token, err := newUseCase.Login(loginWebModel)

			assert.Equal(t, "", token)
			assert.Error(t, errors2.ErrAuthenticationFailed, err)
		}

		t.Logf("\tWhen a failed login authoritzation when user not validate.")
		{
			user := tt.NewRegistrationUser(testName)
			assert.Nil(t, newUseCase.RegisterUser(user))

			loginWebModel := webmodels.LoginUser{
				Username: user.Email,
				Password: user.Password,
			}

			token, err := newUseCase.Login(loginWebModel)

			assert.Equal(t, "", token)
			assert.Error(t, errors2.ErrAuthenticationFailed, err)
		}
	}
}
