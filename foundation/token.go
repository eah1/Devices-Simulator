package foundation

import (
	"github.com/m1/go-generate-password/generator"
	"github.com/pkg/errors"
)

// GenerateToken token with validations.
func GenerateToken(length uint) (*string, error) {
	// Token Generation.
	validationConfig := new(generator.Config)
	validationConfig.Length = length
	validationConfig.IncludeSymbols = false
	validationConfig.IncludeNumbers = true
	validationConfig.IncludeLowercaseLetters = true
	validationConfig.IncludeUppercaseLetters = true
	validationConfig.ExcludeSimilarCharacters = false
	validationConfig.ExcludeAmbiguousCharacters = false

	// Generate token.
	gen, err := generator.New(validationConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Error generating token")
	}

	// Create token.
	token, err := gen.Generate()
	if err != nil {
		return nil, errors.Wrap(err, "generate token error")
	}

	return token, nil
}
