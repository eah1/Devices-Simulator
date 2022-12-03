package v1_test

import (
	"device-simulator/app/services/simulator-api/handlers"
	"device-simulator/business/db/store"
	"device-simulator/business/sys/auth"
	"device-simulator/business/sys/binder"
	"device-simulator/business/web/responses"
	"device-simulator/business/web/webmodels"
	tt "device-simulator/foundation/test"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

const (
	loginURI  = "/api/v1/auth/login"
	logoutURI = "/api/v1/auth/logout"
)

func TestAuthLogin(t *testing.T) {
	t.Parallel()

	testName := "handler-auth-login"

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

	t.Log("Given the need to login auth endpoint.")
	{
		t.Logf("\tWhen we send the body fields with the unwanted format.")
		{
			validator := new(responses.Validator)

			// all fields empty.
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, loginURI, webmodels.LoginUser{}, headers, nil))

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)

			// email format.
			loginAuth := webmodels.LoginUser{
				Username: faker.RandomString(20),
				Password: faker.RandomString(20),
			}

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, loginURI, loginAuth, headers, nil))

			err = json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)

			// password format.
			loginAuth.Password = faker.RandomString(7)

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, loginURI, loginAuth, headers, nil))

			err = json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)
		}

		t.Logf("\tWhen we send a nil body.")
		{
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, loginURI, nil, headers, nil))

			validator := new(responses.Validator)

			err := json.Unmarshal(rec.Body.Bytes(), &validator)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "ERROR", validator.Status)
		}

		t.Logf("\tWhen a correct login.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			loginAuth := webmodels.LoginUser{
				Username: email,
				Password: password,
			}

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, loginURI, loginAuth, headers, nil))

			successLogin := new(responses.SuccessLogin)

			err := json.Unmarshal(rec.Body.Bytes(), &successLogin)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "OK", successLogin.Status)
			assert.GreaterOrEqual(t, len(successLogin.Token), 0)
		}

		t.Logf("\tWhen a failed login when user not exist.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			loginAuth := webmodels.LoginUser{
				Username: faker.Internet().Email(),
				Password: password,
			}

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, loginURI, loginAuth, headers, nil))

			successLogin := new(responses.Failed)

			err := json.Unmarshal(rec.Body.Bytes(), &successLogin)
			require.NoError(t, err)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Equal(t, "ERROR", successLogin.Status)
		}

		t.Logf("\tWhen a failed login when password not correct.")
		{
			email, _ := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			loginAuth := webmodels.LoginUser{
				Username: email,
				Password: faker.RandomString(20),
			}

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, loginURI, loginAuth, headers, nil))

			successLogin := new(responses.Failed)

			err := json.Unmarshal(rec.Body.Bytes(), &successLogin)
			require.NoError(t, err)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Equal(t, "ERROR", successLogin.Status)
		}

		t.Logf("\tWhen a failed login when user is not activate.")
		{
			email, password := tt.RegisterUser(t, app, testName)

			loginAuth := webmodels.LoginUser{
				Username: email,
				Password: password,
			}

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, loginURI, loginAuth, headers, nil))

			successLogin := new(responses.Failed)

			err := json.Unmarshal(rec.Body.Bytes(), &successLogin)
			require.NoError(t, err)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Equal(t, "ERROR", successLogin.Status)
		}
	}
}

func TestAuthLogout(t *testing.T) {
	t.Parallel()

	testName := "handler-auth-logout"

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

	t.Log("Given the need to logout auth endpoint.")
	{
		t.Logf("\tWhen a sending not token.")
		{
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, logoutURI, nil, headers, nil))
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}

		t.Logf("\tWhen a sending format wrong.")
		{
			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer ",
			}

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, logoutURI, nil, headers, nil))
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			headers = map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer ...",
			}

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, logoutURI, nil, headers, nil))
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}

		t.Logf("\tWhen a correct logout.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			token := tt.AuthLogin(t, app, email, password)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + token,
			}

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, logoutURI, nil, headers, nil))
			assert.Equal(t, http.StatusOK, rec.Code)
		}

		t.Logf("\tWhen a failed logout token is not valid.")
		{
			email, password := tt.RegisterUser(t, app, testName)
			tt.ValidationUser(t, app, newStore, email)

			token := tt.AuthLogin(t, app, email, password)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + token,
			}

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, logoutURI, nil, headers, nil))
			assert.Equal(t, http.StatusOK, rec.Code)

			_, rec = tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, logoutURI, nil, headers, nil))
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}

		t.Logf("\tWhen a failed logout witch not exist userId")
		{
			claims := auth.CustomClaims{
				StandardClaims: auth.NewStandardClaims(),
				Email:          faker.Internet().Email(),
				ID:             uuid.NewString(),
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			str, err := token.SignedString([]byte(""))
			require.NoError(t, err)

			headers := map[string]string{
				"Content-Type":  "application/json; charset=utf-8",
				"Authorization": "Bearer " + str,
			}

			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodPost, logoutURI, nil, headers, nil))
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	}
}
