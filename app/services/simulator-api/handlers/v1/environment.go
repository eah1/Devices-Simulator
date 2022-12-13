// Package v1 contains the group v1 and subgroups.
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

type Environment struct {
	cfg     handler.HandlerConfig
	usecase usecase.UseCase
}

// NewEnvironment constructs a handler Environment.
func NewEnvironment(cfg handler.HandlerConfig) Environment {
	return Environment{
		cfg: cfg,
		usecase: usecase.NewUseCase(
			cfg.Log, cfg.Config, store.NewStore(cfg.Log, cfg.DB), cfg.ClientQueue, cfg.EmailSender),
	}
}

// NewEnvironmentServiceGroup create a group environment handlers.
func NewEnvironmentServiceGroup(app *echo.Group, prefix string, handlers Environment) {
	environment := app.Group(prefix)

	handlers.cfg.AuthorizationUser = common.AuthorizationUser(handlers.usecase.GetCore(), handlers.cfg.Log)

	environment.POST("", handlers.Create, handlers.cfg.JWTConfig, handlers.cfg.AuthorizationUser)
}

// Create environment create EndPoint.
// @Summary Create environment EndPoint
// @Tags Environments
// @Description Create a new environment in the system.
// @Param Authorization header string true "Authentication header"
// @Param EnvironmentCreate body webmodels.CreateEnvironment true "EnvironmentCreate"
// @Accept json
// @Produce json
// @Success 201 {object} responses.SuccessEnvironment
// @Failure 400 {object} responses.Validator
// @Failure 401 {object} responses.Failed
// @Failure 500 {object} responses.Failed
// @Router /api/v1/environments [post].
func (h Environment) Create(ctx echo.Context) error {
	user, _ := ctx.Get("user").(models.User)

	createEnvironment := new(webmodels.CreateEnvironment)

	if err := ctx.Bind(createEnvironment); err != nil {
		if strings.Contains(err.Error(), "validator:") {
			return fmt.Errorf("%w", ctx.JSON(http.StatusBadRequest, responses.Validator{
				Status: "ERROR", Details: strings.Split(err.Error()[10:], ","),
			}))
		}

		return fmt.Errorf("%w", ctx.JSON(errors.ErrorHandlingEnvironmentCreate(err, h.cfg.Log)))
	}

	environmentInfo, err := h.usecase.CreateEnvironment(*createEnvironment, user.ID)
	if err != nil {
		return fmt.Errorf("traceid:%s  request.environment.Create: %w", ctx.Request().Header.Get(echo.HeaderXRequestID),
			ctx.JSON(errors.ErrorHandlingEnvironmentCreate(err, h.cfg.Log)))
	}

	return fmt.Errorf("%w", ctx.JSON(http.StatusCreated,
		responses.SuccessEnvironment{Status: "OK", Environment: environmentInfo}))
}
