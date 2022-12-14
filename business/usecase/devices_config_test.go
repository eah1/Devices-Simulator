package usecase_test

import (
	mycDBErrors "device-simulator/business/db/errors"
	"device-simulator/business/db/store"
	mycErrors "device-simulator/business/sys/errors"
	"device-simulator/business/usecase"
	"device-simulator/business/web/webmodels"
	tt "device-simulator/foundation/test"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUseCaseCreateDeviceConfig(t *testing.T) {
	t.Parallel()

	testName := "use-case-create-devices-config"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	newUseCase := usecase.NewUseCase(
		newLog, newConfig, newStore, tt.InitClientQueue(t, newConfig), tt.InitEmailConfig(t, newConfig))

	t.Log("Given the need to work with the create devices config use case.")
	{
		t.Logf("\tWhen a correct a create device config.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			deviceConfigCreate := tt.NewCreateDevicesConfig(testName)

			newDeviceConfig, err := newUseCase.CreateDeviceConfig(deviceConfigCreate, userDB.ID)
			assert.NotEmpty(t, newDeviceConfig)
			assert.Nil(t, err)
			assert.Equal(t, deviceConfigCreate.Name, newDeviceConfig.Name)
		}

		t.Logf("\tWhen a failed create device config when user not exist.")
		{
			deviceConfigCreate := tt.NewCreateDevicesConfig(testName)

			newDeviceConfig, err := newUseCase.CreateDeviceConfig(deviceConfigCreate, uuid.NewString())
			assert.Equal(t, webmodels.InformationDevicesConfig{}, newDeviceConfig)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a failed create device config when data environment is wrong.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			deviceConfigCreate := tt.NewCreateDevicesConfig(testName)
			deviceConfigCreate.Name = "name\000"

			var customError *mycDBErrors.PsqlError

			newDeviceConfig, err := newUseCase.CreateDeviceConfig(deviceConfigCreate, userDB.ID)
			assert.Equal(t, webmodels.InformationDevicesConfig{}, newDeviceConfig)
			assert.Equal(t, true, errors.As(err, &customError))
		}
	}
}
