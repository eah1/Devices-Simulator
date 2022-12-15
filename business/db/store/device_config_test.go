package store_test

import (
	"device-simulator/business/db/store"
	tt "device-simulator/foundation/test"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestDeviceConfigCreate(t *testing.T) {
	t.Parallel()

	testName := "store-device-config-create"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	t.Log("Given the need to work with device config insert database.")
	{
		user := tt.UserCreate(t, newStore, testName)

		t.Logf("\tWhen creating device config.")
		{
			deviceConfig := tt.NewDeviceConfig(testName, user.ID)
			assert.Nil(t, newStore.DeviceConfigCreate(deviceConfig))
		}

		t.Logf("\tWhen creating a duplicate device config.")
		{
			deviceConfig := tt.DeviceConfigCreate(t, newStore, testName, user.ID)
			assert.Error(t, newStore.DeviceConfigCreate(deviceConfig))
		}

		t.Logf("\tWhen creating device config with invalid id.")
		{
			deviceConfig := tt.NewDeviceConfig(testName, user.ID)
			deviceConfig.ID = faker.RandomString(36)
			assert.Error(t, newStore.DeviceConfigCreate(deviceConfig))
		}

		t.Logf("\tWhen creating device config with invalid userId.")
		{
			deviceConfig := tt.NewDeviceConfig(testName, uuid.NewString())
			assert.Error(t, newStore.DeviceConfigCreate(deviceConfig))
		}
	}
}
