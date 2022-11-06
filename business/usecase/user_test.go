package usecase_test

import (
	"testing"

	"device-simulator/business/db/store"
	errors2 "device-simulator/business/sys/errors"
	"device-simulator/business/usecase"
	tt "device-simulator/foundation/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUseCaseRegisterUser(t *testing.T) {
	t.Parallel()

	testName := "use-case-register_user"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newDataBase := tt.InitDatabase(t, newConfig, newLog)
	newStore := store.NewStore(newLog, newDataBase)
	newUseCase := usecase.NewUseCase(newLog, newConfig, newStore)

	t.Log("Given the need to work with the registration user use case.")
	{
		t.Logf("\tWhen a correct a registration user.")
		{
			user := tt.NewRegistrationUser(testName)

			assert.Nil(t, newUseCase.RegisterUser(user))

			query, err := newDataBase.Query("SELECT * FROM users WHERE email = ?", user.Email)
			require.NoError(t, err)

			assert.Equal(t, user.FirstName, string(query[0]["first_name"]))
			assert.Equal(t, user.LastName, string(query[0]["last_name"]))
			assert.Equal(t, user.Email, string(query[0]["email"]))
			assert.Equal(t, user.Company, string(query[0]["company"]))
			assert.Equal(t, user.Language, string(query[0]["language"]))
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
