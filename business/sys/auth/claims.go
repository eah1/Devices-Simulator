// Package auth provides support for authorization configuration.
package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const timeExpired = 24

// CustomClaims are custom claims extending default ones.
type CustomClaims struct {
	jwt.StandardClaims
	Email string
	ID    string
}

// NewStandardClaims crate a standard claim.
func NewStandardClaims() jwt.StandardClaims {
	standardClaims := new(jwt.StandardClaims)
	standardClaims.ExpiresAt = time.Now().Add(time.Hour * timeExpired).Unix()

	return *standardClaims
}
