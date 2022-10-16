// Package logger contains system logger library.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger initialization of logger.
func InitLogger(service, env string) (*zap.SugaredLogger, error) {
	atom := zap.NewAtomicLevel()
	if env == "local" {
		atom.SetLevel(zapcore.ErrorLevel)
	}

	config := zap.NewProductionConfig()
	config.Level = atom
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}
