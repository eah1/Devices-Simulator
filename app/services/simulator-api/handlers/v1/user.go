// Package v1 contains the group v1 and subgroups.
//
//nolint:wrapcheck
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

	handlers.cfg.AuthorizationUser = common.AuthorizationUser(handlers.usecase.GetCore(), handlers.cfg.Log)

	users.POST("", handlers.Create)
	users.POST("/activate/:activateToken", handlers.Activate)
	users.GET("", handlers.Detail, handlers.cfg.JWTConfig, handlers.cfg.AuthorizationUser)
	users.PUT("", handlers.Update, handlers.cfg.JWTConfig, handlers.cfg.AuthorizationUser)
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

		return ctx.JSON(errors.ErrorHandlingUserCreate(err, h.cfg.Log))
	}

	if err := h.usecase.RegisterUser(*userRegister); err != nil {
		return ctx.JSON(errors.ErrorHandlingUserCreate(err, h.cfg.Log))
	}

	if err := h.usecase.SendValidationEmail(userRegister.Email); err != nil {
		return ctx.JSON(errors.ErrorHandlingUserCreate(err, h.cfg.Log))
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
		return ctx.JSON(errors.ErrorHandlingUserActivate(err, h.cfg.Log))
	}

	return ctx.JSON(http.StatusOK, responses.Success{Status: "OK"})
}

// Detail user information EndPoint.
// @Summary Detail user information EndPoint
// @Tags Users
// @Description Detail a user information in the system.
// @Param Authorization header string true "Authentication header"
// @Produce json
// @Success 200 {object} responses.SuccessUser
// @Failure 401 {object} responses.Failed
// @Failure 500 {object} responses.Failed
// @Security ApiKeyAuth
// @Router /api/v1/users [get].
func (h User) Detail(ctx echo.Context) error {
	user, _ := ctx.Get("user").(models.User)

	return ctx.JSON(http.StatusOK, responses.SuccessUser{Status: "OK", User: h.usecase.InformationUser(user)})
}

// Update user information EndPoint.
// @Summary Update user information EndPoint
// @Tags Users
// @Description Update a user information in the system.
// @Param Authorization header string true "Authentication header"
// @Param UserUpdate body webmodels.UpdateUser true "UserUpdate"
// @Accept json
// @Produce json
// @Success 200 {object} responses.Success
// @Failure 400 {object} responses.Validator
// @Failure 401 {object} responses.Failed
// @Failure 500 {object} responses.Failed
// @Security ApiKeyAuth
// @Router /api/v1/users [put].
func (h User) Update(ctx echo.Context) error {
	user, _ := ctx.Get("user").(models.User)

	userUpdate := new(webmodels.UpdateUser)

	if err := ctx.Bind(userUpdate); err != nil {
		if strings.Contains(err.Error(), "validator:") {
			return ctx.JSON(http.StatusBadRequest, responses.Validator{
				Status: "ERROR", Details: strings.Split(err.Error()[10:], ","),
			})
		}

		return ctx.JSON(errors.ErrorHandlingUpdateUser(err, h.cfg.Log))
	}

	if err := h.usecase.UpdateInformationUser(*userUpdate, user.ID); err != nil {
		return fmt.Errorf("%w", ctx.JSON(errors.ErrorHandlingUpdateUser(err, h.cfg.Log)))
	}

	return ctx.JSON(http.StatusOK, responses.Success{Status: "OK"})
}
