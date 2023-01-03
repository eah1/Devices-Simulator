package store_test

import (
	"device-simulator/business/db/store"
	tt "device-simulator/foundation/test"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestDeviceCreate(t *testing.T) {
	t.Parallel()

	testName := "store-device-create"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	t.Log("Given the need to work with device insert database.")
	{
		user := tt.UserCreate(t, newStore, testName)
		environment := tt.EnvironmentCreate(t, newStore, testName, user.ID)
		deviceConfig := tt.DeviceConfigCreate(t, newStore, testName, user.ID)

		t.Logf("\tWhen creating device.")
		{
			device := tt.NewDevice(testName, user.ID, environment.ID, deviceConfig.ID)
			assert.Nil(t, newStore.DeviceCreate(device))
		}

		t.Logf("\tWhen creating a duplicate device.")
		{
			device := tt.DeviceCreate(t, newStore, testName, user.ID, environment.ID, deviceConfig.ID)
			assert.Error(t, newStore.DeviceCreate(device))
		}

		t.Logf("\tWhen creating device with invalid id.")
		{
			device := tt.NewDevice(testName, user.ID, environment.ID, deviceConfig.ID)
			device.ID = faker.RandomString(36)
			assert.Error(t, newStore.DeviceCreate(device))
		}

		t.Logf("\tWhen creating device with invalid userId.")
		{
			device := tt.NewDevice(testName, uuid.NewString(), environment.ID, deviceConfig.ID)
			assert.Error(t, newStore.DeviceCreate(device))
		}

		t.Logf("\tWhen creating device with invalid environmentId.")
		{
			device := tt.NewDevice(testName, user.ID, uuid.NewString(), deviceConfig.ID)
			assert.Error(t, newStore.DeviceCreate(device))
		}

		t.Logf("\tWhen creating device with invalid deviceConfigId.")
		{
			device := tt.NewDevice(testName, user.ID, environment.ID, uuid.NewString())
			assert.Error(t, newStore.DeviceCreate(device))
		}
	}
}
