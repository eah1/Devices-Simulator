package store_test

import (
	mycDBErrors "device-simulator/business/db/errors"
	"device-simulator/business/db/store"
	mycErrors "device-simulator/business/sys/errors"
	tt "device-simulator/foundation/test"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestAuthenticationCreate(t *testing.T) {
	t.Parallel()

	testName := "store-authentication-create"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	t.Log("Given the need to work with authentication insert database.")
	{
		user := tt.UserCreate(t, newStore, testName)

		t.Logf("\tWhen creating authentication.")
		{
			authentication := tt.NewAuthentication(testName, user.ID)
			assert.Nil(t, newStore.AuthenticationCreate(authentication))
		}

		t.Logf("\tWhen creating a duplicate authentication.")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)
			assert.Error(t, newStore.AuthenticationCreate(authentication))
		}

		t.Logf("\tWhen creating authentication with invalid id.")
		{
			authentication := tt.NewAuthentication(testName, user.ID)
			authentication.ID = faker.RandomString(36)
			assert.Error(t, newStore.AuthenticationCreate(authentication))
		}

		t.Logf("\tWhen creating authentication with invalid userId.")
		{
			authentication := tt.NewAuthentication(testName, uuid.NewString())
			assert.Error(t, newStore.AuthenticationCreate(authentication))
		}
	}
}

func TestAuthenticationFindByTokenAndUserID(t *testing.T) {
	t.Parallel()

	testName := "store-authentication-find-by-token-and-userId"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	t.Log("Given the need to work with authentication find by token and userId in database.")
	{
		user := tt.UserCreate(t, newStore, testName)

		t.Logf("\tWhen a correct find authentication by token and userId.")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)

			authenticationFind, err := newStore.AuthenticationFindByTokenAndUserID(authentication.Token, user.ID)
			assert.Equal(t, authentication.ID, authenticationFind.ID)
			assert.Nil(t, err)
		}

		t.Logf("\tWhen a not found find authentication by token and userId which userId not exist.")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)

			authenticationFind, err := newStore.AuthenticationFindByTokenAndUserID(authentication.Token, uuid.NewString())
			assert.Empty(t, authenticationFind)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a not found find authentication by token and userId which token not exist.")
		{
			authenticationFind, err := newStore.AuthenticationFindByTokenAndUserID(faker.RandomString(250), user.ID)
			assert.Empty(t, authenticationFind)
			assert.Error(t, mycErrors.ErrElementNotExist, err)
		}

		t.Logf("\tWhen a not found find authentication by token and userId which format is not correct.")
		{
			authenticationFind, err := newStore.AuthenticationFindByTokenAndUserID("token\000", user.ID)

			var customError *mycDBErrors.PsqlError
			assert.Empty(t, authenticationFind)
			assert.Equal(t, true, errors.As(err, &customError))
		}
	}
}

func TestAuthenticationUpdate(t *testing.T) {
	t.Parallel()

	testName := "store-authentication-update"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))

	t.Log("Given the need to work with user authentication in database.")
	{
		user := tt.UserCreate(t, newStore, testName)

		t.Logf("\tWhen a correct update authentication.")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)
			authentication.LogoutAt = time.Now()
			authentication.Valid = false

			assert.Nil(t, newStore.AuthenticationUpdate(authentication))
		}

		t.Logf("\tWhen a fail authentication update when authentication not exist.")
		{
			authentication := tt.NewAuthentication(testName, user.ID)
			assert.Error(t, mycErrors.ErrElementNotExist, newStore.AuthenticationUpdate(authentication))
		}

		t.Logf("\tWhen a fail authentication update when fields are wrong.")
		{
			authentication := tt.AuthenticationCreate(t, newStore, testName, user.ID)
			authentication.Token = "token\000"

			var customError *mycDBErrors.PsqlError
			assert.Equal(t, true, errors.As(newStore.AuthenticationUpdate(authentication), &customError))
		}
	}
}
