package core_test

import (
	"device-simulator/business/core"
	"device-simulator/business/db/store"
	errors2 "device-simulator/business/sys/errors"
	tt "device-simulator/foundation/test"
	"testing"

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
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

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
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

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

func TestUserCreateValidationToken(t *testing.T) {
	t.Parallel()

	testName := "core-user-create-validation-token"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with a creation validation token.")
	{
		t.Logf("\tWhen a correct crating valiadation token.")
		{
			user := tt.UserCreate(t, newStore, testName)

			assert.Equal(t, "", user.ValidationToken)
			assert.Nil(t, newCore.User.CreateValidationToken(&user))
			assert.Equal(t, 16, len(user.ValidationToken))
		}

		t.Logf("\tWhen a failed creating validation token when user not exist.")
		{
			user := tt.NewUser(testName)

			assert.Equal(t, "", user.ValidationToken)
			assert.Error(t, errors2.ErrElementNotExist, newCore.User.CreateValidationToken(&user))
			assert.Equal(t, "", user.ValidationToken)
		}
	}
}

func TestUserFindByEmail(t *testing.T) {
	t.Parallel()

	testName := "core-user-find-by-email"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with a find user by email.")
	{
		t.Logf("\tWhen a correct find user by email.")
		{
			user := tt.UserCreate(t, newStore, testName)

			userFind, err := newCore.User.FindByEmail(user.Email)

			assert.Equal(t, user.ID, userFind.ID)
			assert.Equal(t, user.Email, userFind.Email)
			assert.Nil(t, err)
		}

		t.Logf("\tWhen a failed find user by email where email not exist.")
		{
			userFind, err := newCore.User.FindByEmail(faker.Internet().Email())

			assert.Empty(t, userFind)
			assert.Error(t, errors2.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a failed find user by email where email wrong format.")
		{
			userFind, err := newCore.User.FindByEmail("")

			assert.Empty(t, userFind)
			assert.Error(t, errors2.ErrElementNotExist, err)

			userFind, err = newCore.User.FindByEmail(faker.RandomString(20))

			assert.Empty(t, userFind)
			assert.Error(t, errors2.ErrElementNotExist, err)
		}
	}
}
