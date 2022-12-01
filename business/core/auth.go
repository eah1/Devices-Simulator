// Package core contains core business API.
package core

import (
	"device-simulator/app/config"
	"device-simulator/business/sys/auth"
	"fmt"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// AuthCore manages the set of API for auth access.
type AuthCore struct {
	log    *zap.SugaredLogger
	config config.Config
	core   *Core
}

// NewAuthCore constructs a core for auth API access.
func NewAuthCore(log *zap.SugaredLogger, config config.Config, core *Core) AuthCore {
	return AuthCore{
		log:    log,
		config: config,
		core:   core,
	}
}

// GenerateToken create token auth.
func (c *AuthCore) GenerateToken(claims auth.CustomClaims) (string, error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	str, err := token.SignedString([]byte(c.config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("core.auth.GenerateToken.SignedString: %w", err)
	}

	return str, nil
}
