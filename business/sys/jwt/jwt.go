// Package jwt provides support for jwt configuration.
package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
)

// NewConfigJWT configure middleware with the custom claims type.
func NewConfigJWT(secretKey string, claims jwt.Claims) *middleware.JWTConfig {
	configJWT := new(middleware.JWTConfig)
	configJWT.Claims = claims
	configJWT.SigningKey = []byte(secretKey)

	return configJWT
}
