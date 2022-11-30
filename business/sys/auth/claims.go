// Package auth provides support for jwt configuration.
package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// CustomClaims are custom claims extending default ones.
type CustomClaims struct {
	jwt.StandardClaims
	Email string
	ID    string
}

// NewStandardClaims crate a standard claim.
func NewStandardClaims() jwt.StandardClaims {
	standardClaims := new(jwt.StandardClaims)
	standardClaims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()

	return *standardClaims
}
