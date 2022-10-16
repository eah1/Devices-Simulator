// Package config content struct Config. The function will parse the environment variables to run the service.
package config

type Config struct {
	Host             string
	HostName         string
	Port             string
	BaseURL          string
	ServerURI        string
	Sentry           string
	Environment      string
	Release          string
	TracesSampleRate float64
}
