package usecase_test

import (
	mycDBErrors "device-simulator/business/db/errors"
	"device-simulator/business/db/store"
	mycErrors "device-simulator/business/sys/errors"
	"device-simulator/business/usecase"
	tt "device-simulator/foundation/test"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestUseCaseRegisterUser(t *testing.T) {
	t.Parallel()

	testName := "use-case-register-user"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newUseCase := usecase.NewUseCase(newLog, newConfig, newStore, nil, nil)

	t.Log("Given the need to work with the registration user use case.")
	{
		t.Logf("\tWhen a correct a registration user.")
		{
			user := tt.NewRegistrationUser(testName)

			assert.Nil(t, newUseCase.RegisterUser(user))

			userDB, err := newStore.UserFindByEmail(user.Email)
			require.NoError(t, err)

			assert.Equal(t, user.FirstName, userDB.FirstName)
			assert.Equal(t, user.LastName, userDB.LastName)
			assert.Equal(t, user.Email, userDB.Email)
			assert.Equal(t, user.Company, userDB.Company)
			assert.Equal(t, user.Language, userDB.Language)
		}

		t.Logf("\tWhen creating a registration duplicate user.")
		{
			user := tt.NewRegistrationUser(testName)

			var customError *mycDBErrors.PsqlError

			assert.Nil(t, newUseCase.RegisterUser(user))
			assert.Equal(t, true, errors.As(newUseCase.RegisterUser(user), &customError))
		}

		t.Logf("\tWhen creating a registration wrong form.")
		{
			user := tt.NewRegistrationUser(testName)
			user.Email = "email\000"

			var customError *mycDBErrors.PsqlError

			assert.Equal(t, true, errors.As(newUseCase.RegisterUser(user), &customError))
		}
	}
}

func TestUseCaseSendValidationEmail(t *testing.T) {
	t.Parallel()

	testName := "use-case-send-validation-token"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	newUseCase := usecase.NewUseCase(
		newLog, newConfig, newStore, tt.InitClientQueue(t, newConfig), tt.InitEmailConfig(t, newConfig))

	t.Log("Given the need to work with the send validation email use case.")
	{
		t.Logf("\tWhen a correct a send validation email.")
		{
			user := tt.NewRegistrationUser(testName)
			assert.Nil(t, newUseCase.RegisterUser(user))

			assert.Nil(t, newUseCase.SendValidationEmail(user.Email))
		}

		t.Logf("\tWhen a failed send validation email where email not exist.")
		{
			assert.Error(t, mycErrors.ErrElementNotExist, newUseCase.SendValidationEmail(faker.Internet().Email()))
		}
	}
}

func TestUseCaseActivateUser(t *testing.T) {
	t.Parallel()

	testName := "use-case-activate-user"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	newUseCase := usecase.NewUseCase(
		newLog, newConfig, newStore, tt.InitClientQueue(t, newConfig), tt.InitEmailConfig(t, newConfig))

	t.Log("Given the need to work with the activate user use case.")
	{
		t.Logf("\tWhen a correct a activate user.")
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
		}

		t.Logf("\tWhen a failed activate user when user not exist.")
		{
			assert.Error(t, mycErrors.ErrElementNotExist, newUseCase.ActivateUser(faker.RandomString(16)))
		}

		t.Logf("\tWhen a failed activate user when user is alredy active.")
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

			assert.Error(t, mycErrors.ErrAuthenticationFailed, newUseCase.ActivateUser(userDB.ValidationToken))
		}
	}
}

func TestUseCaseInformationUser(t *testing.T) {
	t.Parallel()

	testName := "use-case-information-user"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	newUseCase := usecase.NewUseCase(
		newLog, newConfig, newStore, tt.InitClientQueue(t, newConfig), tt.InitEmailConfig(t, newConfig))

	t.Log("Given the need to work with the information user use case.")
	{
		t.Logf("\tWhen a correct a information user.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			userInformationWebModel := newUseCase.InformationUser(userDB)
			assert.Equal(t, userInformationWebModel.ID, userDB.ID)
			assert.Equal(t, userInformationWebModel.FirstName, userDB.FirstName)
			assert.Equal(t, userInformationWebModel.LastName, userDB.LastName)
			assert.Equal(t, userInformationWebModel.Email, userDB.Email)
			assert.Equal(t, userInformationWebModel.Company, userDB.Company)
			assert.Equal(t, userInformationWebModel.Language, userDB.Language)
		}
	}
}

func TestUseCaseUpdateInformationUser(t *testing.T) {
	t.Parallel()

	testName := "use-case-update-user"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	newUseCase := usecase.NewUseCase(
		newLog, newConfig, newStore, tt.InitClientQueue(t, newConfig), tt.InitEmailConfig(t, newConfig))

	t.Log("Given the need to work with the update user use case.")
	{
		t.Logf("\tWhen a correct a update user.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			userUpdate := tt.NewUpdateUser(testName)

			assert.Nil(t, newUseCase.UpdateInformationUser(userUpdate, userDB.ID))
		}

		t.Logf("\tWhen a failed update user when user not exist.")
		{
			userUpdate := tt.NewUpdateUser(testName)

			assert.Error(t, mycErrors.ErrElementNotExist, newUseCase.UpdateInformationUser(userUpdate, uuid.NewString()))
		}

		t.Logf("\tWhen a failed update user when data user is wrong.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			userUpdate := tt.NewUpdateUser(testName)
			userUpdate.FirstName = "fistName\000"

			var customError *mycDBErrors.PsqlError
			assert.Equal(t, true, errors.As(newUseCase.UpdateInformationUser(userUpdate, userDB.ID), &customError))
		}
	}
}

func TestUseCaseUpdatePasswordUser(t *testing.T) {
	t.Parallel()

	testName := "use-case-update-password-user"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	newUseCase := usecase.NewUseCase(
		newLog, newConfig, newStore, tt.InitClientQueue(t, newConfig), tt.InitEmailConfig(t, newConfig))

	t.Log("Given the need to work with the update password user use case.")
	{
		t.Logf("\tWhen a correct a update password user.")
		{
			// Create a register user and validation.
			email, password := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			userUpdatePassword := tt.NewUpdatePasswordUser(password)

			assert.Nil(t, newUseCase.UpdatePasswordUser(userUpdatePassword, userDB.ID))
		}

		t.Logf("\tWhen a failed update password user when user not exist.")
		{
			userUpdatePassword := tt.NewUpdatePasswordUser(faker.RandomString(10))

			assert.Error(t, mycErrors.ErrElementNotExist, newUseCase.UpdatePasswordUser(userUpdatePassword, uuid.NewString()))
		}

		t.Logf("\tWhen a failed update password user when user not correct password.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			userUpdatePassword := tt.NewUpdatePasswordUser(faker.RandomString(10))

			assert.Error(t, mycErrors.ErrAuthenticationFailed, newUseCase.UpdatePasswordUser(userUpdatePassword, userDB.ID))
		}
	}
}
