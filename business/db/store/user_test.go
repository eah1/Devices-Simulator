package store_test

import (
	"device-simulator/business/db/store"
	"device-simulator/business/sys/errors"
	"testing"

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

func TestUserFindByEmail(t *testing.T) {
	t.Parallel()

	testName := "store-user-find-by-email"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	t.Log("Given the need to work with user find by email in database.")
	{
		t.Logf("\tWhen a correct find user by email.")
		{
			user := tt.UserCreate(t, newStore, testName)

			userFind, err := newStore.UserFindByEmail(user.Email)

			assert.Equal(t, user.Email, userFind.Email)
			assert.Nil(t, err)
		}

		t.Logf("\tWhen a not found find user by email which user not exist.")
		{
			userFind, err := newStore.UserFindByEmail(faker.Internet().Email())

			assert.Empty(t, userFind)
			assert.Error(t, errors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a not found find user by email which format email is wrong.")
		{
			userFind, err := newStore.UserFindByEmail("")

			assert.Empty(t, userFind)
			assert.Error(t, errors.ErrElementNotExist, err)

			userFind, err = newStore.UserFindByEmail(faker.RandomString(20))

			assert.Empty(t, userFind)
			assert.Error(t, errors.ErrElementNotExist, err)
		}
	}
}

func TestUserFindByValidationToken(t *testing.T) {
	t.Parallel()

	testName := "store-user-find-by-validation-token"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	t.Log("Given the need to work with user find by validation token in database.")
	{
		t.Logf("\tWhen a correct find user by validation token.")
		{
			user := tt.UserCreate(t, newStore, testName)

			userFind, err := newStore.UserFindByValidationToken(user.ValidationToken)

			assert.Equal(t, user.ValidationToken, userFind.ValidationToken)
			assert.Nil(t, err)
		}

		t.Logf("\tWhen a not found find user by validate token which user not exist.")
		{
			userFind, err := newStore.UserFindByValidationToken(faker.RandomString(16))

			assert.Empty(t, userFind)
			assert.Error(t, errors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a not found find user by validate token which format email is wrong.")
		{
			userFind, err := newStore.UserFindByValidationToken("")

			assert.Empty(t, userFind)
			assert.Error(t, errors.ErrElementNotExist, err)

			userFind, err = newStore.UserFindByValidationToken(faker.RandomString(16))

			assert.Empty(t, userFind)
			assert.Error(t, errors.ErrElementNotExist, err)
		}
	}
}

func TestUserUpdate(t *testing.T) {
	t.Parallel()

	testName := "store-user-update"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	t.Log("Given the need to work with user update in database.")
	{
		t.Logf("\tWhen a correct update user.")
		{
			user := tt.UserCreate(t, newStore, testName)
			user.Language = faker.RandomChoice([]string{"en", "es", "fr", "pt"})

			assert.Nil(t, newStore.UserUpdate(user))
		}

		t.Logf("\tWhen a fail user update when user not exist.")
		{
			user := tt.NewUser(testName)

			assert.Error(t, errors.ErrElementNotExist, newStore.UserUpdate(user))
		}

		t.Logf("\tWhen a fail user update when fields are wrong.")
		{
			user := tt.UserCreate(t, newStore, testName)
			user.Language = faker.RandomChoice([]string{"en", "es", "fr", "pt"})
			user.ID = faker.RandomString(20)

			assert.Error(t, errors.ErrElementRequest, newStore.UserUpdate(user))
		}
	}
}
