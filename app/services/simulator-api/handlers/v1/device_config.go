// Package v1 contains the group v1 and subgroups.
//
//nolint:dupl
package v1

import (
	"device-simulator/business/core/models"
	"device-simulator/business/db/store"
	"device-simulator/business/sys/handler"
	"device-simulator/business/usecase"
	"device-simulator/business/web/errors"
	"device-simulator/business/web/middlewares/common"
	"device-simulator/business/web/responses"
	"device-simulator/business/web/webmodels"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type DeviceConfig struct {
	cfg     handler.HandlerConfig
	usecase usecase.UseCase
}

// NewDeviceConfig constructs a handler DeviceConfig.
func NewDeviceConfig(cfg handler.HandlerConfig) DeviceConfig {
	return DeviceConfig{
		cfg: cfg,
		usecase: usecase.NewUseCase(
			cfg.Log, cfg.Config, store.NewStore(cfg.Log, cfg.DB), cfg.ClientQueue, cfg.EmailSender),
	}
}

// NewDeviceConfigServiceGroup create a group device config handlers.
func NewDeviceConfigServiceGroup(app *echo.Group, prefix string, handlers DeviceConfig) {
	deviceConfig := app.Group(prefix)

	handlers.cfg.AuthorizationUser = common.AuthorizationUser(handlers.usecase.GetCore(), handlers.cfg.Log)

	deviceConfig.POST("", handlers.Create, handlers.cfg.JWTConfig, handlers.cfg.AuthorizationUser)
}

// Create device config create EndPoint.
// @Summary Create device config EndPoint
// @Tags Device-Config
// @Description Create a new device config in the system.
// @Param Authorization header string true "Authentication header"
// @Param DeviceConfigCreate body webmodels.InformationDevicesConfig true "DeviceConfigCreate"
// @Accept json
// @Produce json
// @Success 201 {object} responses.SuccessDeviceConfig
// @Failure 400 {object} responses.Validator
// @Failure 401 {object} responses.Failed
// @Failure 500 {object} responses.Failed
// @Router /api/v1/devices-config [post].
func (h DeviceConfig) Create(ctx echo.Context) error {
	user, _ := ctx.Get("user").(models.User)

	createDeviceConfig := new(webmodels.CreateDeviceConfig)

	if err := ctx.Bind(createDeviceConfig); err != nil {
		if strings.Contains(err.Error(), "validator:") {
			return fmt.Errorf("%w", ctx.JSON(http.StatusBadRequest, responses.Validator{
				Status: "ERROR", Details: strings.Split(err.Error()[10:], ","),
			}))
		}

		return fmt.Errorf("%w", ctx.JSON(errors.ErrorHandlingEnvironmentCreate(err, h.cfg.Log)))
	}

	deviceConfigInfo, err := h.usecase.CreateDeviceConfig(*createDeviceConfig, user.ID)
	if err != nil {
		return fmt.Errorf("traceid:%s  request.device_config.Create: %w", ctx.Request().Header.Get(echo.HeaderXRequestID),
			ctx.JSON(errors.ErrorHandlingEnvironmentCreate(err, h.cfg.Log)))
	}

	return fmt.Errorf("%w", ctx.JSON(http.StatusCreated,
		responses.SuccessDeviceConfig{Status: "OK", DeviceConfig: deviceConfigInfo}))
}
