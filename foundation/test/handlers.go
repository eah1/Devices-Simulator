package test_test

import (
	"device-simulator/business/db/store"
	"device-simulator/business/sys/handler"
	"device-simulator/business/web/responses"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	usersURI         = "/api/v1/users"
	usersActivateURI = "/api/v1/users/activate/"
)

// InitHandlerConfig create a config handler.
func InitHandlerConfig(t *testing.T, service string) handler.HandlerConfig {
	t.Helper()

	handlerConfig := new(handler.HandlerConfig)
	handlerConfig.Log = InitLogger(t, service)
	handlerConfig.Config = InitConfig()
	handlerConfig.DB = InitDatabase(t, handlerConfig.Config, handlerConfig.Log)

	handlerConfig.Config.TemplateFolder = "../../business/template/"

	handlerConfig.EmailSender = InitEmailConfig(t, handlerConfig.Config)
	handlerConfig.ClientQueue = InitClientQueue(t, handlerConfig.Config)

	return *handlerConfig
}

func RegisterUser(t *testing.T, app *echo.Echo, testName string) (string, string) {
	t.Helper()

	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	// Register new user.
	registerUser := NewRegistrationUser(testName)

	_, rec := MakeRequest(t, NewRequestTest(app, http.MethodPost, usersURI, registerUser, headers, nil))

	success := new(responses.Success)

	err := json.Unmarshal(rec.Body.Bytes(), &success)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "OK", success.Status)

	return registerUser.Email, registerUser.Password
}

func ValidationUser(t *testing.T, app *echo.Echo, newStore store.Store, email string) {
	t.Helper()

	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	// find user in database.
	userDB, err := newStore.UserFindByEmail(email)
	require.NoError(t, err)

	// Validate user.
	_, rec := MakeRequest(t, NewRequestTest(app, http.MethodPost,
		usersActivateURI+userDB.ValidationToken, nil, headers, nil))

	success := new(responses.Success)

	err = json.Unmarshal(rec.Body.Bytes(), &success)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", success.Status)
}
