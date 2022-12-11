package binder

import (
	"device-simulator/business/web/validator"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/leebenson/conform"
	"github.com/pkg/errors"
)

type CustomBinder struct{}

func (cb *CustomBinder) Bind(body interface{}, ctx echo.Context) (err error) {
	// You may use default binder
	db := &echo.DefaultBinder{}
	if err = db.Bind(body, ctx); err != nil {
		return fmt.Errorf("binder.Binder: %w", err)
	}

	// format structure.
	if err := conform.Strings(body); err != nil {
		return fmt.Errorf("binder.Conform: %w", err)
	}

	if msg, err := validator.IsValid(body); err != nil {
		return errors.New("validator:" + strings.Join(msg, " , "))
	}

	return nil
}
