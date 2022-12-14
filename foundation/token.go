package foundation

import (
	mycErrors "device-simulator/business/sys/errors"
	"fmt"

	"github.com/m1/go-generate-password/generator"
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
		return nil, fmt.Errorf("token.GenerateToken.New(%d): %v, errMyc: %w", length, err, mycErrors.ErrGenerateToken)
	}

	// Create token.
	token, err := gen.Generate()
	if err != nil {
		return nil, fmt.Errorf("token.GenerateToken.Generate(%d): %v, errMyc: %w", length, err, mycErrors.ErrGenerateToken)
	}

	return token, nil
}
