// Package health contains health handlers.
package health

import (
	"device-simulator/business/sys/handler"
	"device-simulator/business/web/responses"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Health struct {
	cfg handler.HandlerConfig
}

// NewHealth constructs a handler health.
func NewHealth(cfg handler.HandlerConfig) Health {
	return Health{
		cfg: cfg,
	}
}

// Health check EndPoint.
// @Summary Health check EndPoint
// @Tags Health
// @Description Get status service to be alive
// @Accept json
// @Produce json
// @Success 200 {object} responses.SuccessHealth
// @Router /api/health [get].
func (h Health) Health(ctx echo.Context) error {
	respHealth := responses.SuccessHealth{
		Status: "OK", Health: responses.Health{
			BuildVersion: os.Getenv("BUILD_REF"),
		},
	}

	return errors.Wrap(ctx.JSON(http.StatusOK, respHealth), "")
}
