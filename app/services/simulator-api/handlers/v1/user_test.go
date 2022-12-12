package v1_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"device-simulator/app/services/simulator-api/handlers"
	"device-simulator/business/db/store"
	"device-simulator/business/sys/binder"
	"device-simulator/business/web/responses"
	"device-simulator/business/web/webmodels"
	tt "device-simulator/foundation/test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

const (
	usersURI                  = "/api/v1/users"
	usersActivateURI          = "/api/v1/users/activate/"
	usersCredentialsChangeURI = "/api/v1/users/credentials/changePassword"
)

func TestUserCreate(t *testing.T) {
	t.Parallel()

	testName := "handler-user-create"

	// Setup echo.
	app := echo.New()

	// set binder custom.
	app.Binder = &binder.CustomBinder{}

	// Create a configuration handlers.
	handlerConfig := tt.InitHandlerConfig(t, "t-"+testName)

	// Initializing handles.
	handlers.Handlers(app, handlerConfig)

	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	t.Log("Given the need to create user endpoint.")
	{
		t.Logf("\tWhen we send the body fields with the unwanted format.")
		{
			validator := new(responses.Validator)

			// all fields empty.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(
				app, http.MethodPost, usersURI, webmodels.RegisterUser{}, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)

			// email format.
			registerUser := tt.NewRegistrationUser(testName)
			registerUser.Email = faker.RandomString(20)

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersURI, registerUser, headers, nil))

			err = json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)

			// password format.
			registerUser = tt.NewRegistrationUser(testName)
			registerUser.Password = faker.RandomString(7)

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersURI, registerUser, headers, nil))

			err = json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)

			// language format.
			registerUser = tt.NewRegistrationUser(testName)
			registerUser.Language = faker.RandomString(2)

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersURI, registerUser, headers, nil))

			err = json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)
		}

		t.Logf("\tWhen we send a nil body.")
		{
			validator := new(responses.Validator)

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersURI, nil, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)
		}

		t.Logf("\tWhen a correct create a user.")
		{
			success := new(responses.Success)

			registerUser := tt.NewRegistrationUser(testName)

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersURI, registerUser, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &success)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "OK", success.Status)
		}

		t.Logf("\tWhen a duplicate user.")
		{
			success := new(responses.Success)
			failed := new(responses.Failed)

			registerUser := tt.NewRegistrationUser(testName)

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersURI, registerUser, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &success)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "OK", success.Status)

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersURI, registerUser, headers, nil))

			err = json.Unmarshal(rec.Body.Bytes(), &failed)
			require.NoError(t, err)

			assert.Equal(t, http.StatusConflict, rec.Code)
			assert.Equal(t, "ERROR", failed.Status)
		}
	}
}

func TestUserActivate(t *testing.T) {
	t.Parallel()

	testName := "handler-user-validate"

	// Setup echo.
	app := echo.New()

	// set binder custom.
	app.Binder = &binder.CustomBinder{}

	// Create a configuration handlers.
	handlerConfig := tt.InitHandlerConfig(t, "t-"+testName)

	newStore := store.NewStore(handlerConfig.Log, handlerConfig.DB)

	// Initializing handles.
	handlers.Handlers(app, handlerConfig)

	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	t.Log("Given the need to validate user endpoint.")
	{
		t.Logf("\tWhen a correct validate a new user.")
		{
			success := new(responses.Success)

			// Register new user.
			registerUser := tt.NewRegistrationUser(testName)

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersURI, registerUser, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &success)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "OK", success.Status)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(registerUser.Email)
			require.NoError(t, err)

			// Validate user.
			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersActivateURI+userDB.ValidationToken, nil, headers, nil))

			err = json.Unmarshal(rec.Body.Bytes(), &success)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "OK", success.Status)
		}

		t.Logf("\tWhen a validate user not exist.")
		{
			success := new(responses.Success)

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersActivateURI+faker.RandomString(16), nil, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &success)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", success.Status)
		}

		t.Logf("\tWhen a validate user alredy validate.")
		{
			success := new(responses.Success)
			failed := new(responses.Failed)

			// Register new user.
			registerUser := tt.NewRegistrationUser(testName)

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersURI, registerUser, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &success)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "OK", success.Status)

			// find user in database.
			userDB, err := newStore.UserFindByEmail(registerUser.Email)
			require.NoError(t, err)

			// Validate user.
			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersActivateURI+userDB.ValidationToken, nil, headers, nil))

			err = json.Unmarshal(rec.Body.Bytes(), &success)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "OK", success.Status)

			// New validation.
			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, usersActivateURI+userDB.ValidationToken, nil, headers, nil))

			err = json.Unmarshal(rec.Body.Bytes(), &failed)
			require.NoError(t, err)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Equal(t, "ERROR", failed.Status)
		}
	}
}

func TestUserDetail(t *testing.T) {
	t.Parallel()

	testName := "handler-user-details"

	// Setup echo.
	app := echo.New()

	// set binder custom.
	app.Binder = &binder.CustomBinder{}

	// Create a configuration handlers.
	handlerConfig := tt.InitHandlerConfig(t, "t-"+testName)

	newStore := store.NewStore(handlerConfig.Log, handlerConfig.DB)

	// Initializing handles.
	handlers.Handlers(app, handlerConfig)

	t.Log("Given the need to get user details endpoint.")
	{
		t.Logf("\tWhen a correct details of the user.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			successUser := new(responses.SuccessUser)

			// get user details.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodGet, usersURI, nil, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &successUser)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, successUser.User.Email, email)
		}

		t.Logf("\tWhen a sending not token.")
		{
			headers := map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			}

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodGet, usersURI, nil, headers, nil))
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}

		t.Logf("\tWhen a sending format wrong.")
		{
			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer ",
			}

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodGet, usersURI, nil, headers, nil))
			assert.Equal(t, http.StatusUnauthorized, rec.Code)

			headers = map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer ...",
			}

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodGet, usersURI, nil, headers, nil))
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	}
}

func TestUserUpdate(t *testing.T) {
	t.Parallel()

	testName := "handler-user-update"

	// Setup echo.
	app := echo.New()

	// set binder custom.
	app.Binder = &binder.CustomBinder{}

	// Create a configuration handlers.
	handlerConfig := tt.InitHandlerConfig(t, "t-"+testName)

	newStore := store.NewStore(handlerConfig.Log, handlerConfig.DB)

	// Initializing handles.
	handlers.Handlers(app, handlerConfig)

	t.Log("Given the need to get user update endpoint.")
	{
		t.Logf("\tWhen we send the body fields with the unwanted format.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			validator := new(responses.Validator)

			// all fields empty.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPut, usersURI, webmodels.UpdateUser{}, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)

			// language format.
			userUpdate := tt.NewUpdateUser(testName)
			userUpdate.Language = faker.RandomString(2)

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPut, usersURI, userUpdate, headers, nil))

			err = json.Unmarshal(rec.Body.Bytes(), &validator)
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

			validator := new(responses.Validator)

			// update user.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPut, usersURI, nil, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)
		}

		t.Logf("\tWhen a correct update of the user.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			updateUser := tt.NewUpdateUser(testName)

			success := new(responses.Success)

			// update user.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPut, usersURI, updateUser, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &success)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "OK", success.Status)
		}

		t.Logf("\tWhen a failed update a body data wrong.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			updateUser := tt.NewUpdateUser(testName)
			updateUser.FirstName = "fistName\000"

			success := new(responses.Success)

			// update user.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPut, usersURI, updateUser, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &success)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", success.Status)
		}
	}
}

func TestUserChangePassword(t *testing.T) {
	t.Parallel()

	testName := "handler-user-update-password"

	// Setup echo.
	app := echo.New()

	// set binder custom.
	app.Binder = &binder.CustomBinder{}

	// Create a configuration handlers.
	handlerConfig := tt.InitHandlerConfig(t, "t-"+testName)

	newStore := store.NewStore(handlerConfig.Log, handlerConfig.DB)

	// Initializing handles.
	handlers.Handlers(app, handlerConfig)

	t.Log("Given the need to get user update password endpoint.")
	{
		t.Logf("\tWhen we send the body fields with the unwanted format.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			validator := new(responses.Validator)

			// all fields empty.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPut, usersCredentialsChangeURI, webmodels.UpdatePasswordUser{}, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)

			// newPassword equal to currentPassword
			userUpdatePassword := tt.NewUpdatePasswordUser(password)
			userUpdatePassword.NewPassword = password

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPut, usersCredentialsChangeURI, userUpdatePassword, headers, nil))

			err = json.Unmarshal(rec.Body.Bytes(), &validator)
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

			validator := new(responses.Validator)

			// update user.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPut, usersCredentialsChangeURI, nil, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)
		}

		t.Logf("\tWhen a correct update password of the user.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + tt.AuthLogin(t, app, email, password),
			}

			updatePasswordUser := tt.NewUpdatePasswordUser(password)

			success := new(responses.Success)

			// update user.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPut, usersCredentialsChangeURI, updatePasswordUser, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &success)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "OK", success.Status)
		}
	}
}
