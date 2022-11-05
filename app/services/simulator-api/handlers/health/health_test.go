package health_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"device-simulator/app/services/simulator-api/handlers"
	"device-simulator/business/web/responses"
	tt "device-simulator/foundation/test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const healthURI = "/api/health"

func TestHandlerHealth(t *testing.T) {
	t.Parallel()

	// Setup echo.
	app := echo.New()

	// Create a configuration handlers.
	handlerConfig := tt.InitHandlerConfig(t, "test-handler-health")

	// Initializing handles.
	handlers.Handlers(app, handlerConfig)

	t.Log("Given the need to check health endpoint.")
	{
		t.Logf("\tWhen a health makes a request")
		{
			_, rec := tt.MakeRequest(t, tt.NewRequestTest(app, http.MethodGet, healthURI, nil, nil, nil))

			var respData responses.SuccessHealth

			err := json.Unmarshal(rec.Body.Bytes(), &respData)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "OK", respData.Status)
			assert.Equal(t, "", respData.Health.BuildVersion)
		}
	}
}
