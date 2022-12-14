package store_test

import (
	mycDBErrors "device-simulator/business/db/errors"
	"device-simulator/business/db/store"
	mycErrors "device-simulator/business/sys/errors"
	tt "device-simulator/foundation/test"
	"errors"
	"github.com/google/uuid"
	"testing"

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
			assert.Error(t, newStore.UserCreate(user))
		}

		t.Logf("\tWhen creating user with invalid id.")
		{
			user := tt.NewUser(testName)
			user.ID = faker.RandomString(36)
			assert.Error(t, newStore.UserCreate(user))
		}

		t.Logf("\tWhen creating user with invalid insertion table.")
		{
			user := tt.NewUser(testName)
			user.FirstName = "name\000"

			assert.Error(t, newStore.UserCreate(user))
		}
	}
}

func TestUserFindByID(t *testing.T) {
	t.Parallel()

	testName := "store-user-find-by-id"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	t.Log("Given the need to work with user find by id in database.")
	{
		t.Logf("\tWhen a correct find user by id.")
		{
			user := tt.UserCreate(t, newStore, testName)

			userFind, err := newStore.UserFindByID(user.ID)
			assert.Equal(t, user.ID, userFind.ID)
			assert.Nil(t, err)
		}

		t.Logf("\tWhen a not found find user by id which user not exist.")
		{
			userFind, err := newStore.UserFindByID(uuid.NewString())
			assert.Empty(t, userFind)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a not found find user by id which format id is wrong.")
		{
			userFind, err := newStore.UserFindByID("")
			assert.Empty(t, userFind)
			assert.Error(t, mycErrors.ErrElementNotExist, err)

			userFind, err = newStore.UserFindByID(faker.RandomString(20))
			assert.Empty(t, userFind)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a not found find user by id which format is not correct.")
		{
			userFind, err := newStore.UserFindByID("id\000")

			var customError *mycDBErrors.PsqlError
			assert.Empty(t, userFind)
			assert.Equal(t, true, errors.As(err, &customError))
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
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a not found find user by email which format email is wrong.")
		{
			userFind, err := newStore.UserFindByEmail("")
			assert.Empty(t, userFind)
			assert.Error(t, mycErrors.ErrElementNotExist, err)

			userFind, err = newStore.UserFindByEmail(faker.RandomString(20))
			assert.Empty(t, userFind)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a not found find user by email which format is not correct.")
		{
			userFind, err := newStore.UserFindByEmail("email\000")

			var customError *mycDBErrors.PsqlError
			assert.Empty(t, userFind)
			assert.Equal(t, true, errors.As(err, &customError))
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
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a not found find user by validate token which format email is wrong.")
		{
			userFind, err := newStore.UserFindByValidationToken("token\000")

			var customError *mycDBErrors.PsqlError
			assert.Empty(t, userFind)
			assert.Equal(t, true, errors.As(err, &customError))
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
			assert.Error(t, mycErrors.ErrElementNotExist, newStore.UserUpdate(user))
		}

		t.Logf("\tWhen a fail user update when fields are wrong.")
		{
			user := tt.UserCreate(t, newStore, testName)
			user.Language = faker.RandomChoice([]string{"en", "es", "fr", "pt"})
			user.ID = "id\000"

			var customError *mycDBErrors.PsqlError
			assert.Equal(t, true, errors.As(newStore.UserUpdate(user), &customError))
		}
	}
}
