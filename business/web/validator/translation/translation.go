// Package translation contains an error translation errors to text.
package translation

import (
	"device-simulator/business/web/validator/config"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// TranslateEnError translate english error of validator structs.
func TranslateEnError(err error, v *validator.Validate) (errs []string) {
	trans := config.CreateEnTranslation(v)

	return translateError(err, trans)
}

// translateError translate error of validator structs.
func translateError(err error, trans ut.Translator) (errs []string) {
	if err != nil {
		//nolint: errorlint
		validatorErrs, _ := err.(validator.ValidationErrors)

		for _, e := range validatorErrs {
			translatedErr := errors.New(e.Translate(trans))

			errs = append(errs, translatedErr.Error())
		}

		return errs
	}

	return nil
}
