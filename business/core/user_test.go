package core_test

import (
	"testing"

	"device-simulator/business/core"
	"device-simulator/business/db/store"
	errors2 "device-simulator/business/sys/errors"
	tt "device-simulator/foundation/test"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestUserGeneratePassword(t *testing.T) {
	t.Parallel()

	testName := "core-user-generate_password"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore)

	t.Log("Given the need to work with generate password hash.")
	{
		t.Logf("\tWhen a correct geneate password hash.")
		{
			user := tt.NewUser(testName)
			assert.Nil(t, newCore.User.GeneratePassword(faker.Internet().Password(8, 62), &user))
			assert.Nil(t, newCore.User.GeneratePassword("", &user))
		}
	}
}

func TestUserCreate(t *testing.T) {
	t.Parallel()

	testName := "core-user-create"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore)

	t.Log("Given the need to work with a creation user.")
	{
		t.Logf("\tWhen a correct crating user.")
		{
			user := tt.NewUser(testName)

			assert.Nil(t, newCore.User.Create(user))
		}

		t.Logf("\tWhen a create a duplicate user.")
		{
			user := tt.UserCreate(t, newStore, testName)

			assert.Error(t, errors2.ErrElementDuplicated, newCore.User.Create(user))
		}
	}
}
