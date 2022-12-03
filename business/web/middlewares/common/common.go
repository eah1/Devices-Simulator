package common

import (
	"device-simulator/business/core"
	"device-simulator/business/sys/auth"
	"device-simulator/business/web/errors"
	"device-simulator/business/web/responses"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

const timeParse = 1000000

// AddCommonMiddlewares add basics middlewares.
func AddCommonMiddlewares(app *echo.Echo, log *zap.SugaredLogger) {
	app.Use(GenerateTraceID())
	app.Use(ZapLogger(log))
	app.Use(SentryTransaction())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
}

// GenerateTraceID aggregate traceId in HeaderXRequestID.
func GenerateTraceID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Request().Header.Set(echo.HeaderXRequestID, uuid.New().String())

			if err := next(ctx); err != nil {
				ctx.Error(err)
			}

			return nil
		}
	}
}

// ZapLogger aggregate logger in request.
func ZapLogger(log *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()
			req, res := ctx.Request(), ctx.Response()

			requestID := getRequestID(ctx)

			log.Infow("request started", "traceid", requestID, "method", req.Method,
				"path", req.URL.Path, "remoteaddr", req.RemoteAddr)

			if err := next(ctx); err != nil {
				ctx.Error(err)
			}

			log.Infow("request completed", "traceid", requestID, "method", req.Method,
				"path", req.URL.Path, "remoteaddr", req.RemoteAddr, "statuscode", res.Status,
				"exec_time", float64(time.Since(start))/timeParse)

			return nil
		}
	}
}

// getRequestId get request id.
func getRequestID(ctx echo.Context) string {
	if ctx.Request().Header.Get(echo.HeaderXRequestID) != "" {
		return ctx.Request().Header.Get(echo.HeaderXRequestID)
	}

	if ctx.Response().Header().Get(echo.HeaderXRequestID) != "" {
		return ctx.Response().Header().Get(echo.HeaderXRequestID)
	}

	return "00000000-0000-0000-0000-000000000000"
}

// SentryTransaction send transaction into sentry.
func SentryTransaction() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			span := sentry.StartSpan(
				ctx.Request().Context(),
				ctx.Request().RequestURI,
				sentry.TransactionName(fmt.Sprintf(ctx.Request().RequestURI)))

			span.StartTime = time.Now()

			if err := next(ctx); err != nil {
				ctx.Error(err)
			}

			span.EndTime = time.Now()
			span.Status = sentry.SpanStatus(ctx.Response().Status)
			span.Finish()

			return nil
		}
	}
}

// AuthorizationUser check authorization user.
func AuthorizationUser(core core.Core, log *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			token, ok := ctx.Get("user").(*jwt.Token)
			if !ok || token == nil {
				return next(ctx)
			}

			claims, ok := token.Claims.(*auth.CustomClaims)
			if !ok {
				return fmt.Errorf("%w", ctx.JSON(http.StatusUnauthorized,
					responses.Failed{Status: "ERROR", Error: "error in get token claims"}))
			}

			authentication, err := core.Authentication.FindByTokenAndUserID(token.Raw, claims.ID)
			if err != nil {
				return fmt.Errorf("%w", ctx.JSON(errors.ErrorHandlingLogout(err, log)))
			}

			if err := core.Authentication.IsValid(authentication); err != nil {
				return fmt.Errorf("%w", ctx.JSON(errors.ErrorHandlingLogout(err, log)))
			}

			user, err := core.User.FindByEmail(claims.Email)
			if err != nil {
				return fmt.Errorf("%w", ctx.JSON(errors.ErrorHandlingLogout(err, log)))
			}

			ctx.Set("token", token)
			ctx.Set("user", user)

			return next(ctx)
		}
	}
}
