package usecase_test

import (
	mycDBErrors "device-simulator/business/db/errors"
	"device-simulator/business/db/store"
	mycErrors "device-simulator/business/sys/errors"
	"device-simulator/business/usecase"
	"device-simulator/business/web/webmodels"
	tt "device-simulator/foundation/test"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUseCaseCreateDevice(t *testing.T) {
	t.Parallel()

	testName := "use-case-create-devices"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	newUseCase := usecase.NewUseCase(
		newLog, newConfig, newStore, tt.InitClientQueue(t, newConfig), tt.InitEmailConfig(t, newConfig))

	t.Log("Given the need to work with the create device use case.")
	{
		t.Logf("\tWhen a correct a create device.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			newEnvironment, err := newUseCase.CreateEnvironment(tt.NewCreateEnvironment(testName), userDB.ID)
			assert.NotEmpty(t, newEnvironment)
			assert.Nil(t, err)

			newDeviceConfig, err := newUseCase.CreateDeviceConfig(tt.NewCreateDevicesConfig(testName), userDB.ID)
			assert.NotEmpty(t, newDeviceConfig)
			assert.Nil(t, err)

			deviceCreate := tt.NewCreateDevice(testName, newEnvironment.ID, newDeviceConfig.ID)
			newDevice, err := newUseCase.CreateDevice(deviceCreate, userDB.ID)
			assert.NotEmpty(t, newDevice)
			assert.Nil(t, err)
			assert.Equal(t, deviceCreate.Name, newDevice.Name)
		}

		t.Logf("\tWhen a failed create device when user not exist.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			newEnvironment, err := newUseCase.CreateEnvironment(tt.NewCreateEnvironment(testName), userDB.ID)
			assert.NotEmpty(t, newEnvironment)
			assert.Nil(t, err)

			newDeviceConfig, err := newUseCase.CreateDeviceConfig(tt.NewCreateDevicesConfig(testName), userDB.ID)
			assert.NotEmpty(t, newDeviceConfig)
			assert.Nil(t, err)

			deviceCreate := tt.NewCreateDevice(testName, newEnvironment.ID, newDeviceConfig.ID)
			newDevice, err := newUseCase.CreateDevice(deviceCreate, uuid.NewString())
			assert.Equal(t, webmodels.InformationDevice{}, newDevice)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a failed create device when environment not exist.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			newEnvironment, err := newUseCase.CreateEnvironment(tt.NewCreateEnvironment(testName), userDB.ID)
			assert.NotEmpty(t, newEnvironment)
			assert.Nil(t, err)

			newDeviceConfig, err := newUseCase.CreateDeviceConfig(tt.NewCreateDevicesConfig(testName), userDB.ID)
			assert.NotEmpty(t, newDeviceConfig)
			assert.Nil(t, err)

			deviceCreate := tt.NewCreateDevice(testName, uuid.NewString(), newDeviceConfig.ID)
			newDevice, err := newUseCase.CreateDevice(deviceCreate, uuid.NewString())
			assert.Equal(t, webmodels.InformationDevice{}, newDevice)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a failed create device when device config not exist.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			newEnvironment, err := newUseCase.CreateEnvironment(tt.NewCreateEnvironment(testName), userDB.ID)
			assert.NotEmpty(t, newEnvironment)
			assert.Nil(t, err)

			newDeviceConfig, err := newUseCase.CreateDeviceConfig(tt.NewCreateDevicesConfig(testName), userDB.ID)
			assert.NotEmpty(t, newDeviceConfig)
			assert.Nil(t, err)

			deviceCreate := tt.NewCreateDevice(testName, newEnvironment.ID, uuid.NewString())
			newDevice, err := newUseCase.CreateDevice(deviceCreate, uuid.NewString())
			assert.Equal(t, webmodels.InformationDevice{}, newDevice)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a failed create device when data is wrong.")
		{
			// Create a register user and validation.
			email, _ := tt.UseCaseRegisterValidate(t, newUseCase, newStore, testName)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(email)
			require.NoError(t, err)

			newEnvironment, err := newUseCase.CreateEnvironment(tt.NewCreateEnvironment(testName), userDB.ID)
			assert.NotEmpty(t, newEnvironment)
			assert.Nil(t, err)

			newDeviceConfig, err := newUseCase.CreateDeviceConfig(tt.NewCreateDevicesConfig(testName), userDB.ID)
			assert.NotEmpty(t, newDeviceConfig)
			assert.Nil(t, err)

			deviceCreate := tt.NewCreateDevice(testName, newEnvironment.ID, newDeviceConfig.ID)
			deviceCreate.Name = "deviceName\000"

			var customError *mycDBErrors.PsqlError

			newDevice, err := newUseCase.CreateDevice(deviceCreate, userDB.ID)
			assert.Equal(t, webmodels.InformationDevice{}, newDevice)
			assert.Equal(t, true, errors.As(err, &customError))
		}
	}
}
