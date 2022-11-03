// Package v1 contains the group v1 and subgroups.
package v1

import (
	"device-simulator/business/sys/handler"

	"github.com/labstack/echo/v4"
)

type User struct {
	cfg handler.HandlerConfig
}

// NewUser constructs a handler User.
func NewUser(cfg handler.HandlerConfig) User {
	return User{
		cfg: cfg,
	}
}

// NewUserServiceGroup create a group user handlers.
func NewUserServiceGroup(app *echo.Group, prefix string, handlers User) {
	users := app.Group(prefix)

	users.POST("", handlers.Create)
}

// Create user create EndPoint.
func (h User) Create(ctx echo.Context) error {
	return nil
}
