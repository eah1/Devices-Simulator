// Package validator validation structs web models.
package validator

import (
	"device-simulator/business/web/validator/translation"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// IsValid validation data structures.
func IsValid(dataStruct interface{}) ([]string, error) {
	validate := validator.New()

	if err := validate.Struct(dataStruct); err != nil {
		messageError := translation.TranslateEnError(err, validate)

		return messageError, fmt.Errorf("[%+v]: %w", &dataStruct, err)
	}

	return nil, nil
}
