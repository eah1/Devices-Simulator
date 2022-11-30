// Package sentry contains sentry library connector.
package sentry

import (
	"strings"

	"device-simulator/app/config"
	"github.com/getsentry/sentry-go"
)

// InitSentryConfig initialization of sentry client options.
func InitSentryConfig(config config.Config) sentry.ClientOptions {
	sentryClient := new(sentry.ClientOptions)

	sentryClient.Dsn = config.Sentry
	sentryClient.Environment = config.Environment
	sentryClient.Release = config.Release
	sentryClient.TracesSampler = sentry.TracesSamplerFunc(func(ctx sentry.SamplingContext) sentry.Sampled {
		if strings.Contains(ctx.Span.Op, "/health") {
			return sentry.SampledFalse
		}

		return sentry.UniformTracesSampler(config.TracesSampleRate).Sample(ctx)
	})

	return *sentryClient
}
