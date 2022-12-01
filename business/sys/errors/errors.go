// Package errors contains a error handling.
package errors

import (
	"github.com/pkg/errors"
)

var (
	ErrDB = errors.New("Database error")

	ErrElementDuplicated    = errors.New("Element duplicated")
	ErrElementRequest       = errors.New("Element request failed")
	ErrElementNotExist      = errors.New("Element not exist")
	ErrAuthenticationFailed = errors.New("Authentication failed")
)
