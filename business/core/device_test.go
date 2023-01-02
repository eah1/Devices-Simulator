package core_test

import (
	"device-simulator/business/core"
	mycDBErrors "device-simulator/business/db/errors"
	"device-simulator/business/db/store"
	tt "device-simulator/foundation/test"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeviceCreate(t *testing.T) {
	t.Parallel()

	testName := "core-device-create"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with a creation device.")
	{
		user := tt.UserCreate(t, newStore, testName)
		environment := tt.EnvironmentCreate(t, newStore, testName, user.ID)
		deviceConfig := tt.DeviceConfigCreate(t, newStore, testName, user.ID)

		t.Logf("\tWhen a correct creating device.")
		{
			device := tt.NewDevice(testName, user.ID, environment.ID, deviceConfig.ID)
			assert.Nil(t, newCore.Device.Create(&device))
		}

		t.Logf("\tWhen a create a device user not exist.")
		{
			device := tt.NewDevice(testName, user.ID, environment.ID, deviceConfig.ID)
			device.UserID = uuid.NewString()

			var customError *mycDBErrors.PsqlError
			assert.Equal(t, true, errors.As(newCore.Device.Create(&device), &customError))
		}

		t.Logf("\tWhen a create a device environment not exist.")
		{
			device := tt.NewDevice(testName, user.ID, environment.ID, deviceConfig.ID)
			device.EnvironmentID = uuid.NewString()

			var customError *mycDBErrors.PsqlError
			assert.Equal(t, true, errors.As(newCore.Device.Create(&device), &customError))
		}

		t.Logf("\tWhen a create a device device config not exist.")
		{
			device := tt.NewDevice(testName, user.ID, environment.ID, deviceConfig.ID)
			device.DeviceConfigID = uuid.NewString()

			var customError *mycDBErrors.PsqlError
			assert.Equal(t, true, errors.As(newCore.Device.Create(&device), &customError))
		}

		t.Logf("\tWhen a create a device data wrong.")
		{
			device := tt.NewDevice(testName, user.ID, environment.ID, deviceConfig.ID)
			device.Name = "deviceName\000"

			var customError *mycDBErrors.PsqlError
			assert.Equal(t, true, errors.As(newCore.Device.Create(&device), &customError))
		}
	}
}
