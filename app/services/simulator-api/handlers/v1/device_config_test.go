package v1_test

import (
	"device-simulator/app/services/simulator-api/handlers"
	"device-simulator/business/db/store"
	"device-simulator/business/sys/binder"
	"device-simulator/business/web/responses"
	"device-simulator/business/web/webmodels"
	tt "device-simulator/foundation/test"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const deviceConfigURI = "/api/v1/devices-config"

func TestDeviceConfigCreate(t *testing.T) {
	t.Parallel()

	testName := "handler-device-config-create"

	// Setup echo.
	app := echo.New()

	// set binder custom.
	app.Binder = &binder.CustomBinder{}

	// Create a configuration handlers.
	handlerConfig := tt.InitHandlerConfig(t, "t-"+testName)

	newStore := store.NewStore(handlerConfig.Log, handlerConfig.DB)

	// Initializing handles.
	handlers.Handlers(app, handlerConfig)

	t.Log("Given the need to create device-config endpoint.")
	{
		t.Logf("\tWhen we send the body fields with the unwanted format.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			// all fields empty.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(
				app, http.MethodPost, deviceConfigURI, webmodels.CreateDeviceConfig{}, headers, nil))

			validator := new(responses.Validator)

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)
		}

		t.Logf("\tWhen we send a nil body.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, deviceConfigURI, nil, headers, nil))

			validator := new(responses.Validator)

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)
		}

		t.Logf("\tWhen a correct create device config.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			deviceConfigCreate := tt.NewCreateDevicesConfig(testName)

			// Created device config.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, deviceConfigURI, deviceConfigCreate, headers, nil))

			successDeviceConfig := new(responses.SuccessDeviceConfig)

			err := json.Unmarshal(rec.Body.Bytes(), &successDeviceConfig)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "OK", successDeviceConfig.Status)
			assert.Equal(t, deviceConfigCreate.Name, successDeviceConfig.DeviceConfig.Name)
		}

		t.Logf("\tWhen a faild create device config when data is wrong format.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			deviceConfigCreate := tt.NewCreateDevicesConfig(testName)
			deviceConfigCreate.Name = "name\000"

			// Created device config.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, deviceConfigURI, deviceConfigCreate, headers, nil))

			successFailed := new(responses.Failed)

			err := json.Unmarshal(rec.Body.Bytes(), &successFailed)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", successFailed.Status)
		}
	}
}
