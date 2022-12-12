package store_test

import (
	"device-simulator/business/db/store"
	tt "device-simulator/foundation/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
	"testing"
)

func TestEnvironmentCreate(t *testing.T) {
	t.Parallel()

	testName := "store-environment-create"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	t.Log("Given the need to work with environment insert database.")
	{
		user := tt.UserCreate(t, newStore, testName)

		t.Logf("\tWhen creating environment.")
		{
			environment := tt.NewEnvironment(testName, user.ID)
			assert.Nil(t, newStore.EnvironmentCreate(environment))
		}

		t.Logf("\tWhen creating a duplicate environment.")
		{
			environment := tt.EnvironmentCreate(t, newStore, testName, user.ID)
			assert.Error(t, newStore.EnvironmentCreate(environment))
		}

		t.Logf("\tWhen creating environment with invalid id.")
		{
			environment := tt.NewEnvironment(testName, user.ID)
			environment.ID = faker.RandomString(36)
			assert.Error(t, newStore.EnvironmentCreate(environment))
		}

		t.Logf("\tWhen creating environment with invalid userId.")
		{
			environment := tt.NewEnvironment(testName, uuid.NewString())
			assert.Error(t, newStore.EnvironmentCreate(environment))
		}
	}
}
