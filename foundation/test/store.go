package test_test

import (
	"device-simulator/business/core/models"
	"device-simulator/business/db/store"
	"testing"

	"github.com/stretchr/testify/require"
)

// UserCreate create a user in database.
func UserCreate(t *testing.T, store store.Store, testName string) models.User {
	t.Helper()

	user := NewUser(testName)

	err := store.UserCreate(user)
	require.NoError(t, err)

	return user
}
