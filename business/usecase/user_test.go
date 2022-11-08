package usecase_test

import (
	"device-simulator/business/db/store"
	errors2 "device-simulator/business/sys/errors"
	"device-simulator/business/usecase"
	tt "device-simulator/foundation/test"
	"testing"

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

			assert.Nil(t, newUseCase.RegisterUser(user))
			assert.Error(t, errors2.ErrElementDuplicated, newUseCase.RegisterUser(user))
		}

		t.Logf("\tWhen creating a registration wrong form.")
		{
			user := tt.NewRegistrationUser(testName)
			user.Email = ""

			assert.Error(t, errors2.ErrElementRequest, newUseCase.RegisterUser(user))
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
			assert.Error(t, errors2.ErrElementNotExist, newUseCase.SendValidationEmail(faker.Internet().Email()))
		}
	}
}
