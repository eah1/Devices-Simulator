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
		cfg:     cfg,
		usecase: usecase.NewUseCase(cfg.Log, cfg.Config, store.NewStore(cfg.Log, cfg.DB)),
	}
}

// NewUserServiceGroup create a group user handlers.
func NewUserServiceGroup(app *echo.Group, prefix string, handlers User) {
	users := app.Group(prefix)

	users.POST("", handlers.Create)
}

// Create user create EndPoint.
// @Summary Create Account EndPoint
// @Tags Users
// @Description Create a new user.
// @Param UserRegister body webmodels.RegisterUser true "UserRegister"
// @Accept json
// @Produce json
// @Success 200 {object} responseV2.AccountCreateSuccess
// @Failure 400 {object} response.Validator
// @Failure 409 {object} response.Failed
// @Failure 500 {object} response.Failed.
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
		return err
	}

	return nil
}
