package core_test

import (
	"device-simulator/business/core"
	mycDBErrors "device-simulator/business/db/errors"
	"device-simulator/business/db/store"
	mycErrors "device-simulator/business/sys/errors"
	tt "device-simulator/foundation/test"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestAuthenticationCreate(t *testing.T) {
	t.Parallel()

	testName := "core-authentication-create"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with a creation authentication.")
	{
		t.Logf("\tWhen a correct creating authentication")
		{
			user := tt.UserCreate(t, newStore, testName)

			t.Logf("\tWhen a correct crating authentication.")
			{
				authentication := tt.NewAuthentication(testName, user.ID)
				assert.Nil(t, newCore.Authentication.Create(authentication))
			}

			t.Logf("\tWhen a create a duplicate authentication.")
			{
				authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)

				var customError *mycDBErrors.PsqlError
				assert.Equal(t, true, errors.As(newCore.Authentication.Create(authentication), &customError))
			}
		}
	}
}

func TestAuthenticationFindByTokenAndUserID(t *testing.T) {
	t.Parallel()

	testName := "core-authentication-find-by-token-and-userId"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with a find authentication by token and userId.")
	{
		user := tt.UserCreate(t, newStore, testName)

		t.Logf("\tWhen a correct find authentication by token and userId.")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)

			authenticationFind, err := newCore.Authentication.FindByTokenAndUserID(authentication.Token, user.ID)
			assert.Equal(t, authentication.ID, authenticationFind.ID)
			assert.Equal(t, user.ID, authenticationFind.UserID)
			assert.Nil(t, err)
		}

		t.Logf("\tWhen a failed find authentication by token and userId when token not exist.")
		{
			authenticationFind, err := newCore.Authentication.FindByTokenAndUserID(faker.RandomString(200), user.ID)
			assert.Empty(t, authenticationFind)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a failed find authentication by token and userId when userId not exist.")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)

			authenticationFind, err := newCore.Authentication.FindByTokenAndUserID(authentication.Token, uuid.NewString())
			assert.Empty(t, authenticationFind)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("When a failed find authentication by token and userId when token wrong format")
		{
			authenticationFind, err := newCore.Authentication.FindByTokenAndUserID("token\000", user.ID)

			var customError *mycDBErrors.PsqlError
			assert.Empty(t, authenticationFind)
			assert.Equal(t, true, errors.As(err, &customError))
		}
	}
}

func TestAuthenticationIsValid(t *testing.T) {
	t.Parallel()

	testName := "core-authentication-is-valid"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with a validate authentication.")
	{
		user := tt.UserCreate(t, newStore, testName)

		t.Logf("\tWhen a correct validate authentication.")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)
			assert.Nil(t, newCore.Authentication.IsValid(authentication))
		}

		t.Logf("\tWhen a failed when authentication is not valid")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)
			authentication.Valid = false

			err := newStore.AuthenticationUpdate(authentication)
			assert.Nil(t, err)

			assert.Error(t, mycErrors.ErrAuthenticationFailed, newCore.Authentication.IsValid(authentication))
		}
	}
}

func TestAuthenticationInvalidation(t *testing.T) {
	t.Parallel()

	testName := "core-authentication-invalidation"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with a validate authentication.")
	{
		user := tt.UserCreate(t, newStore, testName)

		t.Logf("\tWhen a correct invalidation authentication.")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)

			assert.True(t, authentication.Valid)
			assert.Nil(t, newCore.Authentication.Invalidation(&authentication))
			assert.False(t, authentication.Valid)
		}

		t.Logf("\tWhen a failed invalidation authentication when authentication is not valid.")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)
			authentication.Valid = false

			assert.Nil(t, newStore.AuthenticationUpdate(authentication))
			assert.False(t, authentication.Valid)

			assert.Error(t, mycErrors.ErrAuthenticationFailed, newCore.Authentication.Invalidation(&authentication))
		}

		t.Logf("\tWhen a failed invalidation authentication when autentication not exist.")
		{
			authentication := tt.NewAuthentication(testName, user.ID)
			assert.Error(t, mycErrors.ErrElementNotExist, newCore.Authentication.Invalidation(&authentication))
		}

		t.Logf("\tWhen a failed invalidation authentication when userId not exist.")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)
			authentication.UserID = uuid.NewString()

			var customError *mycDBErrors.PsqlError
			assert.True(t, authentication.Valid)
			assert.Equal(t, true, errors.As(newCore.Authentication.Invalidation(&authentication), &customError))
			assert.True(t, authentication.Valid)
		}
	}
}
