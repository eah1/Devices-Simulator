// Package errors contains a error handling.
package errors

import (
	"github.com/pkg/errors"
)

var (
	ErrPsqlDefault = errors.New("Default postgresql error")

	ErrElementDuplicated    = errors.New("Element duplicated")
	ErrElementRequest       = errors.New("Element request failed")
	ErrElementNotExist      = errors.New("Element not exist")
	ErrAuthenticationFailed = errors.New("Authentication failed")
)
