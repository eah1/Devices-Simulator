// Package errors contains the logic returns errors.
package errors

import (
	mycDBErrors "device-simulator/business/db/errors"
	mycErrors "device-simulator/business/sys/errors"
	"device-simulator/business/web/responses"
	"errors"
	"net/http"

	sentryGo "github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

// ErrorHandlingUserCreate error handling user create codes.
func ErrorHandlingUserCreate(err error, log *zap.SugaredLogger) (int, responses.Failed) {
	log.Error(err)

	var customError *mycDBErrors.PsqlError

	switch {
	case errors.Is(err, mycErrors.ErrElementNotExist):
		return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
	case errors.Is(err, mycErrors.ErrGenerateToken):
		return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
	case errors.As(err, &customError):
		switch customError.CodeSQL {
		case "23505":
			return http.StatusConflict, responses.Failed{Status: "ERROR", Error: "Email has already exist in the system"}
		default:
			return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
		}
	default:
		sentryGo.CaptureException(err)

		return http.StatusInternalServerError, responses.Failed{Status: "ERROR", Error: "Internal server error"}
	}
}

// ErrorHandlingUserActivate error handling user activate codes.
func ErrorHandlingUserActivate(err error, log *zap.SugaredLogger) (int, responses.Failed) {
	log.Error(err)

	switch {
	case errors.Is(err, mycErrors.ErrElementNotExist):
		return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
	case errors.Is(err, mycErrors.ErrAuthenticationFailed):
		return http.StatusUnauthorized, responses.Failed{Status: "ERROR", Error: "Authentication failed"}
	default:
		sentryGo.CaptureException(err)

		return http.StatusInternalServerError, responses.Failed{Status: "ERROR", Error: "Internal server error"}
	}
}

// ErrorHandlingLogin error handling login codes.
func ErrorHandlingLogin(err error, log *zap.SugaredLogger) (int, responses.Failed) {
	log.Error(err)

	switch {
	case errors.Is(err, mycErrors.ErrElementNotExist) ||
		errors.Is(err, mycErrors.ErrAuthenticationFailed):
		return http.StatusUnauthorized, responses.Failed{Status: "ERROR", Error: "Authentication failed"}
	default:
		sentryGo.CaptureException(err)

		return http.StatusInternalServerError, responses.Failed{Status: "ERROR", Error: "Internal server error"}
	}
}

// ErrorHandlingLogout error handling authorization user codes.
func ErrorHandlingLogout(err error, log *zap.SugaredLogger) (int, responses.Failed) {
	log.Error(err)

	var customError *mycDBErrors.PsqlError

	switch {
	case errors.Is(err, mycErrors.ErrElementNotExist) ||
		errors.Is(err, mycErrors.ErrAuthenticationFailed):
		return http.StatusUnauthorized, responses.Failed{Status: "ERROR", Error: "Authentication failed"}
	case errors.As(err, &customError):
		return http.StatusUnauthorized, responses.Failed{Status: "ERROR", Error: "Authentication failed"}
	default:
		sentryGo.CaptureException(err)

		return http.StatusInternalServerError, responses.Failed{Status: "ERROR", Error: "Internal server error"}
	}
}

// ErrorHandlingUpdateUser error handling user update codes.
func ErrorHandlingUpdateUser(err error, log *zap.SugaredLogger) (int, responses.Failed) {
	log.Error(err)

	var customError *mycDBErrors.PsqlError

	switch {
	case errors.Is(err, mycErrors.ErrElementNotExist):
		return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
	case errors.As(err, &customError):
		switch customError.CodeSQL {
		case "22021":
			return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
		default:
			return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
		}
	default:
		sentryGo.CaptureException(err)

		return http.StatusInternalServerError, responses.Failed{Status: "ERROR", Error: "Internal server error"}
	}
}

// ErrorHandlingUpdatePasswordUser error handling user update password codes.
func ErrorHandlingUpdatePasswordUser(err error, log *zap.SugaredLogger) (int, responses.Failed) {
	log.Error(err)

	var customError *mycDBErrors.PsqlError

	switch {
	case errors.Is(err, mycErrors.ErrElementNotExist) || errors.Is(err, mycErrors.ErrAuthenticationFailed):
		return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
	case errors.Is(err, mycErrors.ErrAuthenticationFailed):
		return http.StatusUnauthorized, responses.Failed{Status: "ERROR", Error: "Authentication failed"}
	case errors.As(err, &customError):
		switch customError.CodeSQL {
		case "22021":
			return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
		default:
			return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
		}
	default:
		sentryGo.CaptureException(err)

		return http.StatusInternalServerError, responses.Failed{Status: "ERROR", Error: "Internal server error"}
	}
}

// ErrorHandlingEnvironmentCreate error handling environment create code.
func ErrorHandlingEnvironmentCreate(err error, log *zap.SugaredLogger) (int, responses.Failed) {
	log.Error(err)

	var customError *mycDBErrors.PsqlError

	switch {
	case errors.Is(err, mycErrors.ErrElementNotExist):
		return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
	case errors.Is(err, mycErrors.ErrAuthenticationFailed):
		return http.StatusUnauthorized, responses.Failed{Status: "ERROR", Error: "Authentication failed"}
	case errors.As(err, &customError):
		switch customError.CodeSQL {
		case "22021":
			return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
		default:
			return http.StatusBadRequest, responses.Failed{Status: "ERROR", Error: "Request failed"}
		}
	default:
		sentryGo.CaptureException(err)

		return http.StatusInternalServerError, responses.Failed{Status: "ERROR", Error: "Internal server error"}
	}
}
