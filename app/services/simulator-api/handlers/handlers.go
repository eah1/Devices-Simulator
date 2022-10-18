// Package handlers contains handlers from the services.
package handlers

import (
	_ "device-simulator/app/services/simulator-api/docs"
	"device-simulator/app/services/simulator-api/handlers/health"
	"device-simulator/business/sys/handler"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Handlers constructs application routes defined.
// @title Swagger MYC-DEVICE-SIMULATOR API
// @version 1.0
// @description Devices Simulator documentation API.
// @BasePath /
// @in header
// @name Authorization.
func Handlers(app *echo.Echo, cfg handler.HandlerConfig) {
	root := app.Group("/api/")

	// Initialize swagger.
	wrapHandler := echoSwagger.WrapHandler
	root.GET("swagger/*", wrapHandler)

	GroupRoot(root, cfg)
}

// GroupRoot create routes in app server.
func GroupRoot(root *echo.Group, cfg handler.HandlerConfig) {
	// Initialize health controllers.
	handlerHealth := health.NewHealth(cfg)
	root.GET("health", handlerHealth.Health)
}
