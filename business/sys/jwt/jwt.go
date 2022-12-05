// Package jwt provides support for jwt configuration.
package jwt

import (
	"device-simulator/business/web/responses"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewConfigJWT configure middleware with the custom claims type.
func NewConfigJWT(secretKey string, claims jwt.Claims) *middleware.JWTConfig {
	configJWT := new(middleware.JWTConfig)
	configJWT.Claims = claims
	configJWT.SigningKey = []byte(secretKey)
	configJWT.ErrorHandlerWithContext = func(err error, ctx echo.Context) error {
		return fmt.Errorf("%w", ctx.JSON(http.StatusUnauthorized,
			responses.Failed{Status: "ERROR", Error: "invalid or expired jwt"}))
	}

	return configJWT
}
