// Package v1 contains the group v1 and subgroups.
//
//nolint:wrapcheck
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

type Auth struct {
	cfg     handler.HandlerConfig
	usecase usecase.UseCase
}

// NewAuth constructs a handler Auth.
func NewAuth(cfg handler.HandlerConfig) Auth {
	return Auth{
		cfg: cfg,
		usecase: usecase.NewUseCase(
			cfg.Log, cfg.Config, store.NewStore(cfg.Log, cfg.DB), cfg.ClientQueue, cfg.EmailSender),
	}
}

// NewAuthServiceGroup create a group auth handlers.
func NewAuthServiceGroup(app *echo.Group, prefix string, handlers Auth) {
	auth := app.Group(prefix)

	auth.POST("/login", handlers.Login)
}

// Login user EndPoint.
// @Summary Login user EndPoint
// @Tags Auth
// @Description Login authentication user in platform
// @Accept json
// @Produce json
// @Param User body webmodels.LoginUser true "LoginUser"
// @Success 200 {object} responses.SuccessLogin
// @Failure 400 {object} responses.Validator
// @Failure 401 {object} responses.Failed
// @Failure 404 {object} responses.Failed
// @Failure 500 {object} responses.Failed
// @Router /api/v1/auth/login [post].
func (h Auth) Login(ctx echo.Context) error {
	userLogin := new(webmodels.LoginUser)

	if err := ctx.Bind(userLogin); err != nil {
		if strings.Contains(err.Error(), "validator:") {
			return ctx.JSON(http.StatusBadRequest, responses.Validator{
				Status: "ERROR", Details: strings.Split(err.Error()[10:], ","),
			})
		}

		return ctx.JSON(errors.ErrorHandlingLogin(err, h.cfg.Log))
	}

	token, err := h.usecase.Login(*userLogin)
	if err != nil {
		return ctx.JSON(errors.ErrorHandlingLogin(err, h.cfg.Log))
	}

	return ctx.JSON(http.StatusOK, responses.SuccessLogin{Status: "OK", Token: token})
}
