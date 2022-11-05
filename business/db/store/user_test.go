package store_test

import (
	"testing"

	"device-simulator/business/db/store"
	"device-simulator/business/sys/errors"
	tt "device-simulator/foundation/test"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestUserCreate(t *testing.T) {
	t.Parallel()

	testName := "store-user-create"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	t.Log("Given the need to work with user insert database.")
	{
		t.Logf("\tWhen creating user.")
		{
			user := tt.NewUser(testName)

			assert.Nil(t, newStore.UserCreate(user))
		}

		t.Logf("\tWhen creating a duplicate user.")
		{
			user := tt.UserCreate(t, newStore, testName)

			assert.Error(t, errors.ErrElementDuplicated, newStore.UserCreate(user))
		}

		t.Logf("\tWhen creating user with invalid id.")
		{
			user := tt.NewUser(testName)
			user.ID = faker.RandomString(36)

			assert.Error(t, errors.ErrElementRequest, newStore.UserCreate(user))
		}
	}
}
