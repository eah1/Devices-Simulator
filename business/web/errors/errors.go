// Package errors contains the logic returns errors.
package errors

import (
	"net/http"

	"device-simulator/business/web/responses"
	"go.uber.org/zap"
)

// HandlingError error handling codes.
func HandlingError(err error, log *zap.SugaredLogger) (int, responses.Failed) {
	return http.StatusInternalServerError, responses.Failed{Status: "ERROR", Error: "Internal server error"}
}
