// Package v1 contains the group v1 and subgroups.
package v1

import (
	"device-simulator/business/db/store"
	"device-simulator/business/sys/handler"
	"device-simulator/business/usecase"
	"device-simulator/business/web/errors"
	"device-simulator/business/web/responses"
	"device-simulator/business/web/webmodels"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type User struct {
	cfg     handler.HandlerConfig
	usecase usecase.UseCase
}

// NewUser constructs a handler User.
func NewUser(cfg handler.HandlerConfig) User {
	return User{
		cfg: cfg,
		usecase: usecase.NewUseCase(
			cfg.Log, cfg.Config, store.NewStore(cfg.Log, cfg.DB), cfg.ClientQueue, cfg.EmailSender),
	}
}

// NewUserServiceGroup create a group user handlers.
func NewUserServiceGroup(app *echo.Group, prefix string, handlers User) {
	users := app.Group(prefix)

	users.POST("", handlers.Create)
	users.POST("/activate/:activateToken", handlers.Activate)
}

// Create user create EndPoint.
// @Summary Create user registration EndPoint
// @Tags Users
// @Description Create a new user in the system.
// @Param UserRegister body webmodels.RegisterUser true "UserRegister"
// @Accept json
// @Produce json
// @Success 201 {object} responses.Success
// @Failure 400 {object} responses.Validator
// @Failure 409 {object} responses.Failed
// @Failure 500 {object} responses.Failed
// @Router /api/v1/users [post].
func (h User) Create(ctx echo.Context) error {
	userRegister := new(webmodels.RegisterUser)

	if err := ctx.Bind(userRegister); err != nil {
		if strings.Contains(err.Error(), "validator:") {
			return ctx.JSON(http.StatusBadRequest, responses.Validator{
				Status: "ERROR", Details: strings.Split(err.Error()[10:], ","),
			})
		}

		h.cfg.Log.Errorw("Body error in bind", "service", "USER", "error", err.Error())

		return ctx.JSON(errors.HandlingError(err, h.cfg.Log))
	}

	if err := h.usecase.RegisterUser(*userRegister); err != nil {
		h.cfg.Log.Errorw("Create -> RegisterUser",
			"service", "HANDLER | USER CREATE | USE CASE USER", "error", err.Error())

		return ctx.JSON(errors.HandlingError(err, h.cfg.Log))
	}

	if err := h.usecase.SendValidationEmail(userRegister.Email); err != nil {
		h.cfg.Log.Errorw("Create -> SendValidationEmail",
			"service", "HANDLER | USER CREATE | USE CASE USER", "error", err.Error())

		return ctx.JSON(errors.HandlingError(err, h.cfg.Log))
	}

	return ctx.JSON(http.StatusCreated, responses.Success{Status: "OK"})
}

// Activate user activation EndPoint.
// @Summary Activate user activation EndPoint
// @Tags Users
// @Description Activation a user in the system.
// @Param activateToken   path      string  true  "ActivateToken"
// @Produce json
// @Success 200 {object} responses.Success
// @Failure 400 {object} responses.Validator
// @Failure 401 {object} responses.Failed
// @Failure 500 {object} responses.Failed
// @Router /api/v1/users/activate/{activateToken} [post].
func (h User) Activate(ctx echo.Context) error {
	activateToken := ctx.Param("activateToken")

	if err := h.usecase.ActivateUser(activateToken); err != nil {
		h.cfg.Log.Errorw("Activate -> ActivateUser",
			"service", "HANDLER | USER CREATE | USE CASE USER", "error", err.Error())

		return ctx.JSON(errors.HandlingError(err, h.cfg.Log))
	}

	return ctx.JSON(http.StatusOK, responses.Success{Status: "OK"})
}
