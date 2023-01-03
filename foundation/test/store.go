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

// AuthenticationCreate create a authentication in database.
func AuthenticationCreate(t *testing.T, store store.Store, testName, userID string) models.Authentication {
	t.Helper()

	authentication := NewAuthentication(testName, userID)

	err := store.AuthenticationCreate(authentication)
	require.NoError(t, err)

	return authentication
}

// EnvironmentCreate create a environment in database.
func EnvironmentCreate(t *testing.T, store store.Store, testName, userID string) models.Environment {
	t.Helper()

	environment := NewEnvironment(testName, userID)

	err := store.EnvironmentCreate(environment)
	require.NoError(t, err)

	return environment
}

// DeviceConfigCreate create a device config in database.
func DeviceConfigCreate(t *testing.T, store store.Store, testName, userID string) models.DeviceConfig {
	t.Helper()

	deviceConfig := NewDeviceConfig(testName, userID)

	err := store.DeviceConfigCreate(deviceConfig)
	require.NoError(t, err)

	return deviceConfig
}

// DeviceCreate create a device in database.
func DeviceCreate(
	t *testing.T, store store.Store, testName, userID, environmentID, deviceConfigID string,
) models.Device {
	t.Helper()

	device := NewDevice(testName, userID, environmentID, deviceConfigID)

	err := store.DeviceCreate(device)
	require.NoError(t, err)

	return device
}
