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

const environmentURI = "/api/v1/environments"

func TestEnvironmentCreate(t *testing.T) {
	t.Parallel()

	testName := "handler-environment-create"

	// Setup echo.
	app := echo.New()

	// set binder custom.
	app.Binder = &binder.CustomBinder{}

	// Create a configuration handlers.
	handlerConfig := tt.InitHandlerConfig(t, "t-"+testName)

	newStore := store.NewStore(handlerConfig.Log, handlerConfig.DB)

	// Initializing handles.
	handlers.Handlers(app, handlerConfig)

	t.Log("Given the need to create environment endpoint.")
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
				app, http.MethodPost, environmentURI, webmodels.CreateEnvironment{}, headers, nil))

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

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, environmentURI, nil, headers, nil))

			validator := new(responses.Validator)

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)
		}

		t.Logf("\tWhen a correct create environment.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			environmentCreate := tt.NewCreateEnvironment(testName)

			// Created environment.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, environmentURI, environmentCreate, headers, nil))

			successEnvironment := new(responses.SuccessEnvironment)

			err := json.Unmarshal(rec.Body.Bytes(), &successEnvironment)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "OK", successEnvironment.Status)
			assert.Equal(t, environmentCreate.Name, successEnvironment.Environment.Name)
		}

		t.Logf("\tWhen a faild create environment when data is wrong format.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			environmentCreate := tt.NewCreateEnvironment(testName)
			environmentCreate.Name = "name\000"

			// Created environment.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, environmentURI, environmentCreate, headers, nil))

			successEnvironment := new(responses.SuccessEnvironment)

			err := json.Unmarshal(rec.Body.Bytes(), &successEnvironment)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", successEnvironment.Status)
		}
	}
}
