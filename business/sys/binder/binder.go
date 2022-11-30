package binder

import (
	"strings"

	"device-simulator/business/web/validator"
	"github.com/labstack/echo/v4"
	"github.com/leebenson/conform"
	"github.com/pkg/errors"
)

type CustomBinder struct{}

func (cb *CustomBinder) Bind(body interface{}, ctx echo.Context) (err error) {
	// You may use default binder
	db := &echo.DefaultBinder{}
	if err = db.Bind(body, ctx); err != nil {
		return err
	}

	// format structure.
	if err := conform.Strings(body); err != nil {
		return errors.New("conform:" + err.Error())
	}

	if msg, err := validator.IsValid(body); err != nil {
		return errors.New("validator:" + strings.Join(msg, " , "))
	}

	return nil
}
