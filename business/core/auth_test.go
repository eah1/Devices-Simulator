package core_test

import (
	"device-simulator/business/core"
	"device-simulator/business/db/store"
	tt "device-simulator/foundation/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthGenerateToken(t *testing.T) {
	t.Parallel()

	testName := "core-auth-generate-token"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()
	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, nil)

	t.Log("Given the need to work with generate token.")
	{
		t.Logf("\tWhen a correct geneate token.")
		{
			claims := tt.NewCustomClaims(testName)

			token, err := newCore.Auth.GenerateToken(claims)
			assert.GreaterOrEqual(t, len(token), 0)
			assert.Nil(t, err)
		}
	}
}
