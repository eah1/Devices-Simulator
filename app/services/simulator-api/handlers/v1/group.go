// Package v1 contains the group v1 and subgroups.
package v1

import (
	"device-simulator/business/sys/handler"

	"github.com/labstack/echo/v4"
)

// CreateGroup create group v1.
func CreateGroup(app *echo.Group, cfg handler.HandlerConfig) {
	v1 := app.Group("v1")

	NewUserServiceGroup(v1, "/users", NewUser(cfg))
}
