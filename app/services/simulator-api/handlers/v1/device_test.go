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

const deviceURI = "/api/v1/devices"

func TestDeviceCreate(t *testing.T) {
	t.Parallel()

	testName := "handler-device-create"

	// Setup echo.
	app := echo.New()

	// set binder custom.
	app.Binder = &binder.CustomBinder{}

	// Create a configuration handlers.
	handlerConfig := tt.InitHandlerConfig(t, "t-"+testName)

	newStore := store.NewStore(handlerConfig.Log, handlerConfig.DB)

	// Initializing handles.
	handlers.Handlers(app, handlerConfig)

	t.Log("Given the need to create device endpoint.")
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
				app, http.MethodPost, deviceURI, webmodels.CreateDevice{}, headers, nil))

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

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, deviceURI, nil, headers, nil))

			validator := new(responses.Validator)

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)
		}

		t.Logf("\tWhen a correct create device.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			// Created environment.
			environmentCreate := tt.NewCreateEnvironment(testName)

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, environmentURI, environmentCreate, headers, nil))

			successEnvironment := new(responses.SuccessEnvironment)

			err := json.Unmarshal(rec.Body.Bytes(), &successEnvironment)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "OK", successEnvironment.Status)

			// Created device config.
			deviceConfigCreate := tt.NewCreateDevicesConfig(testName)

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, deviceConfigURI, deviceConfigCreate, headers, nil))

			successDeviceConfig := new(responses.SuccessDeviceConfig)

			err = json.Unmarshal(rec.Body.Bytes(), &successDeviceConfig)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "OK", successDeviceConfig.Status)

			// Created device.
			deviceCreate := tt.NewCreateDevice(testName, successEnvironment.Environment.ID, successDeviceConfig.DeviceConfig.ID)

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, deviceURI, deviceCreate, headers, nil))

			successDevice := new(responses.SuccessDevice)

			err = json.Unmarshal(rec.Body.Bytes(), &successDevice)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "OK", successDevice.Status)
			assert.Equal(t, deviceCreate.Name, successDevice.Device.Name)
		}

		t.Logf("\tWhen a faild create device when data is wrong format.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			// Created environment.
			environmentCreate := tt.NewCreateEnvironment(testName)

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, environmentURI, environmentCreate, headers, nil))

			successEnvironment := new(responses.SuccessEnvironment)

			err := json.Unmarshal(rec.Body.Bytes(), &successEnvironment)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "OK", successEnvironment.Status)

			// Created device config.
			deviceConfigCreate := tt.NewCreateDevicesConfig(testName)

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, deviceConfigURI, deviceConfigCreate, headers, nil))

			successDeviceConfig := new(responses.SuccessDeviceConfig)

			err = json.Unmarshal(rec.Body.Bytes(), &successDeviceConfig)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "OK", successDeviceConfig.Status)

			// Created device.
			deviceCreate := tt.NewCreateDevice(testName, successEnvironment.Environment.ID, successDeviceConfig.DeviceConfig.ID)
			deviceCreate.Name = "deviceName\000"

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, deviceURI, deviceCreate, headers, nil))

			successFailed := new(responses.Failed)

			err = json.Unmarshal(rec.Body.Bytes(), &successFailed)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", successFailed.Status)
		}
	}
}
