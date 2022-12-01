package core_test

import (
	"device-simulator/business/core"
	mycDBErrors "device-simulator/business/db/errors"
	"device-simulator/business/db/store"
	mycErrors "device-simulator/business/sys/errors"
	tt "device-simulator/foundation/test"
	"errors"
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

func TestUserCheckCredentials(t *testing.T) {
	t.Parallel()

	testName := "core-user_check-credentials"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with check credentials.")
	{
		t.Logf("\tWhen a correct check password.")
		{
			user := tt.UserCreate(t, newStore, testName)
			assert.Nil(t, newCore.User.Activate(&user))
			assert.Nil(t, newCore.User.IsActivate(user))
		}

		t.Logf("\tWhen a wrong check password.")
		{
			user := tt.UserCreate(t, newStore, testName)
			assert.Error(t, mycErrors.ErrAuthenticationFailed, newCore.User.IsActivate(user))
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

			var customError *mycDBErrors.PsqlError

			assert.Equal(t, true, errors.As(newCore.User.Create(user), &customError))
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

			assert.Nil(t, newCore.User.CreateValidationToken(&user))
			assert.Equal(t, 16, len(user.ValidationToken))
		}

		t.Logf("\tWhen a failed creating validation token when user not exist.")
		{
			user := tt.NewUser(testName)

			assert.Error(t, mycErrors.ErrElementNotExist, newCore.User.CreateValidationToken(&user))
			assert.Equal(t, "", user.ValidationToken)
		}
	}
}

func TestUserActivate(t *testing.T) {
	t.Parallel()

	testName := "core-user-activate"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with a activate user.")
	{
		t.Logf("\tWhen a correct validate user.")
		{
			user := tt.UserCreate(t, newStore, testName)

			assert.False(t, user.Validated)
			assert.Nil(t, newCore.User.Activate(&user))
			assert.True(t, user.Validated)
		}

		t.Logf("\tWhen a failed validate user when user not exist.")
		{
			user := tt.NewUser(testName)

			assert.False(t, user.Validated)
			assert.Error(t, mycErrors.ErrElementNotExist, newCore.User.Activate(&user))
			assert.False(t, user.Validated)
		}

		t.Logf("\tWhen a failed when user is activate")
		{
			user := tt.UserCreate(t, newStore, testName)

			assert.False(t, user.Validated)
			assert.Nil(t, newCore.User.Activate(&user))
			assert.True(t, user.Validated)

			assert.Error(t, mycErrors.ErrAuthenticationFailed, newCore.User.Activate(&user))
		}
	}
}

func TestUserIsActivate(t *testing.T) {
	t.Parallel()

	testName := "core-user-is_activate"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work check is validate a user.")
	{
		t.Logf("\tWhen a correct validate.")
		{
			user := tt.NewUser(testName)
			assert.Nil(t, newCore.User.CheckCredentials(user, "password"))
		}

		t.Logf("\tWhen a wrong validate.")
		{
			user := tt.NewUser(testName)
			assert.Error(t, mycErrors.ErrAuthenticationFailed, newCore.User.CheckCredentials(user, faker.RandomString(20)))
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

		t.Logf("\tWhen a failed find user by email when email not exist.")
		{
			userFind, err := newCore.User.FindByEmail(faker.Internet().Email())

			assert.Empty(t, userFind)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a failed find user by email when email wrong format.")
		{
			userFind, err := newCore.User.FindByEmail("email\000")

			var customError *mycDBErrors.PsqlError

			assert.Empty(t, userFind)
			assert.Equal(t, true, errors.As(err, &customError))
		}
	}
}

func TestFindByValidationToken(t *testing.T) {
	t.Parallel()

	testName := "core-user-find-by-validation-token"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with a find user by validation token.")
	{
		t.Logf("\tWhen a correct find user by validation token.")
		{
			user := tt.UserCreate(t, newStore, testName)

			userFind, err := newCore.User.FindByValidationToken(user.ValidationToken)

			assert.Equal(t, user.ID, userFind.ID)
			assert.Equal(t, user.ValidationToken, userFind.ValidationToken)
			assert.Nil(t, err)
		}

		t.Logf("\tWhen a failed find user by validation token when token not exist.")
		{
			userFind, err := newCore.User.FindByValidationToken(faker.RandomString(16))

			assert.Empty(t, userFind)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a failed find user by validation token when token wrong format.")
		{
			userFind, err := newCore.User.FindByValidationToken("token\000")

			var customError *mycDBErrors.PsqlError

			assert.Empty(t, userFind)
			assert.Equal(t, true, errors.As(err, &customError))
		}
	}
}
