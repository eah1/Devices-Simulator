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

func TestDeviceConfigCreate(t *testing.T) {
	t.Parallel()

	testName := "core-device-config-create"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with a creation device config.")
	{
		user := tt.UserCreate(t, newStore, testName)

		t.Logf("\tWhen a correct creating device config.")
		{
			deviceConfig := tt.NewDeviceConfig(testName, user.ID)
			assert.Nil(t, newCore.DeviceConfig.Create(&deviceConfig))
		}

		t.Logf("\tWhen a create a device config user not exist.")
		{
			deviceConfig := tt.NewDeviceConfig(testName, user.ID)
			deviceConfig.UserID = uuid.NewString()

			var customError *mycDBErrors.PsqlError
			assert.Equal(t, true, errors.As(newCore.DeviceConfig.Create(&deviceConfig), &customError))
		}

		t.Logf("\tWhen a create a device config data wrong.")
		{
			deviceConfig := tt.NewDeviceConfig(testName, user.ID)
			deviceConfig.Name = "name\000"

			var customError *mycDBErrors.PsqlError
			assert.Equal(t, true, errors.As(newCore.DeviceConfig.Create(&deviceConfig), &customError))
		}
	}
}
