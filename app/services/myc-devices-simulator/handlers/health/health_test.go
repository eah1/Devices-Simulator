package health_test

import (
	"device-simulator/app/services/myc-devices-simulator/handlers"
	"device-simulator/business/web/responses"
	test_test "device-simulator/foundation/test"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHandlerHealth(t *testing.T) {
	t.Parallel()

	// Setup echo.
	app := echo.New()

	// Create a configuration handlers.
	handlerConfig := test_test.InitHandlerConfig(t, "test-handler-health")

	// Initializing handles.
	handlers.Handlers(app, handlerConfig)

	t.Log("Given the need to check health endpoint.")
	{
		t.Logf("\tWhen a health makes a request")
		{
			_, rec := test_test.MakeRequest(t, test_test.NewRequestTest(app, "GET", "/api/health", nil, nil, nil))

			var respData responses.SuccessHealth

			err := json.Unmarshal(rec.Body.Bytes(), &respData)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "OK", respData.Status)
			assert.Equal(t, "", respData.Health.BuildVersion)
		}
	}
}
