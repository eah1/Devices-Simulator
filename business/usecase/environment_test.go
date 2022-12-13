package usecase_test

import (
	"device-simulator/business/web/webmodels"
	"errors"
	"testing"

	mycDBErrors "device-simulator/business/db/errors"
	"device-simulator/business/db/store"
	mycErrors "device-simulator/business/sys/errors"
	"device-simulator/business/usecase"
	tt "device-simulator/foundation/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUseCaseCreateEnvironment(t *testing.T) {
	t.Parallel()

	testName := "use-case-create-environment"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	newUseCase := usecase.NewUseCase(
		newLog, newConfig, newStore, tt.InitClientQueue(t, newConfig), tt.InitEmailConfig(t, newConfig))

	t.Log("Given the need to work with the create environment use case.")
	{
		t.Logf("\tWhen a correct a create environment.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			environmentCreate := tt.NewCreateEnvironment(testName)

			newEnvironment, err := newUseCase.CreateEnvironment(environmentCreate, userDB.ID)
			assert.NotEmpty(t, newEnvironment)
			assert.Nil(t, err)
		}

		t.Logf("\tWhen a failed create environment when user not exist.")
		{
			environmentCreate := tt.NewCreateEnvironment(testName)

			newEnvironment, err := newUseCase.CreateEnvironment(environmentCreate, uuid.NewString())
			assert.Equal(t, webmodels.InformationEnvironment{}, newEnvironment)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a failed create environment when data environment is wrong.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			environmentCreate := tt.NewCreateEnvironment(testName)
			environmentCreate.Name = "name\000"

			var customError *mycDBErrors.PsqlError

			newEnvironment, err := newUseCase.CreateEnvironment(environmentCreate, userDB.ID)
			assert.Equal(t, webmodels.InformationEnvironment{}, newEnvironment)
			assert.Equal(t, true, errors.As(err, &customError))
		}
	}
}
