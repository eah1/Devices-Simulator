// Package errors contains the logic returns errors.
package errors

import (
	"device-simulator/business/web/responses"
	"net/http"

	errorsMyc "device-simulator/business/sys/errors"

	sentryGo "github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

// HandlingError error handling codes.
func HandlingError(err error, log *zap.SugaredLogger) (int, responses.Failed) {
	switch err {
	case errorsMyc.ErrElementDuplicated:
		return http.StatusConflict, responses.Failed{Status: "ERROR", Error: err.Error()}
	case errorsMyc.ErrElementRequest:
		return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: err.Error()}
	case errorsMyc.ErrAuthenticationFailed:
		return http.StatusUnauthorized, responses.Failed{Status: "ERROR", Error: err.Error()}
	default:
		log.Error(err)
		sentryGo.CaptureException(err)

		return http.StatusInternalServerError, responses.Failed{Status: "ERROR", Error: "Internal server error"}
	}
}
