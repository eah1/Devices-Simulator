package core_test

import (
	"errors"
	"testing"

	"device-simulator/business/core"
	mycDBErrors "device-simulator/business/db/errors"
	"device-simulator/business/db/store"
	tt "device-simulator/foundation/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentCreate(t *testing.T) {
	t.Parallel()

	testName := "core-environment-create"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with a creation environment.")
	{
		user := tt.UserCreate(t, newStore, testName)

		t.Logf("\tWhen a correct creating environment.")
		{
			environment := tt.NewEnvironment(testName, user.ID)
			assert.Nil(t, newCore.Environment.Create(&environment))
		}

		t.Logf("\tWhen a create a environment user not exist.")
		{
			environment := tt.NewEnvironment(testName, user.ID)
			environment.UserID = uuid.NewString()

			var customError *mycDBErrors.PsqlError
			assert.Equal(t, true, errors.As(newCore.Environment.Create(&environment), &customError))
		}

		t.Logf("\tWhen a create a environment data wrong.")
		{
			environment := tt.NewEnvironment(testName, user.ID)
			environment.Name = "name\000"

			var customError *mycDBErrors.PsqlError
			assert.Equal(t, true, errors.As(newCore.Environment.Create(&environment), &customError))
		}
	}
}
