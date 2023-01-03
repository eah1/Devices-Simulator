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

type Device struct {
	cfg     handler.HandlerConfig
	usecase usecase.UseCase
}

// NewDevice constructs a handler Device.
func NewDevice(cfg handler.HandlerConfig) Device {
	return Device{
		cfg: cfg,
		usecase: usecase.NewUseCase(
			cfg.Log, cfg.Config, store.NewStore(cfg.Log, cfg.DB), cfg.ClientQueue, cfg.EmailSender),
	}
}

// NewDeviceServiceGroup create a group device handlers.
func NewDeviceServiceGroup(app *echo.Group, prefix string, handlers Device) {
	device := app.Group(prefix)

	handlers.cfg.AuthorizationUser = common.AuthorizationUser(handlers.usecase.GetCore(), handlers.cfg.Log)

	device.POST("", handlers.Create, handlers.cfg.JWTConfig, handlers.cfg.AuthorizationUser)
}

// Create device create EndPoint.
// @Summary Create device EndPoint
// @Tags Device
// @Description Create a new device in the system.
// @Param Authorization header string true "Authentication header"
// @Param DeviceCreate body webmodels.CreateDevice true "DeviceCreate"
// @Accept json
// @Produce json
// @Success 201 {object} responses.SuccessDevice
// @Failure 400 {object} responses.Validator
// @Failure 401 {object} responses.Failed
// @Failure 500 {object} responses.Failed
// @Router /api/v1/devices [post].
func (h Device) Create(ctx echo.Context) error {
	user, _ := ctx.Get("user").(models.User)

	createDevice := new(webmodels.CreateDevice)

	if err := ctx.Bind(createDevice); err != nil {
		if strings.Contains(err.Error(), "validator:") {
			return fmt.Errorf("%w", ctx.JSON(http.StatusBadRequest, responses.Validator{
				Status: "ERROR", Details: strings.Split(err.Error()[10:], ","),
			}))
		}

		return fmt.Errorf("%w", ctx.JSON(errors.ErrorHandlingEnvironmentCreate(err, h.cfg.Log)))
	}

	deviceInfo, err := h.usecase.CreateDevice(*createDevice, user.ID)
	if err != nil {
		return fmt.Errorf("traceid:%s  request.device.Create: %w", ctx.Request().Header.Get(echo.HeaderXRequestID),
			ctx.JSON(errors.ErrorHandlingEnvironmentCreate(err, h.cfg.Log)))
	}

	return fmt.Errorf("%w", ctx.JSON(http.StatusCreated,
		responses.SuccessDevice{Status: "OK", Device: deviceInfo}))
}
