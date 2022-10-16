package test_test

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RequestTest struct {
	app      *echo.Echo
	method   string
	endpoint string
	data     interface{}
	header   map[string]string
	query    map[string]string
}

func NewRequestTest(app *echo.Echo, method, endpoint string, data interface{},
	headers, queries map[string]string,
) RequestTest {
	return RequestTest{app, method, endpoint, data, headers, queries}
}

func MakeRequest(t *testing.T, requestTest RequestTest) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()

	var payload io.Reader = nil

	if requestTest.data != nil {
		body, err := json.Marshal(requestTest.data)
		require.NoError(t, err)

		payload = bytes.NewBuffer(body)
	}

	request := httptest.NewRequest(requestTest.method, requestTest.endpoint, payload)

	for key, value := range requestTest.header {
		request.Header.Set(key, value)
	}

	if requestTest.query != nil {
		query := request.URL.Query()

		for key, value := range requestTest.query {
			query.Add(key, value)
		}

		request.URL.RawQuery = query.Encode()
	}

	response := httptest.NewRecorder()

	requestTest.app.ServeHTTP(response, request)

	return request, response
}
