// Package config content struct Config. The function will parse the environment variables to run the service.
package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoadConfig loads environment vars into Config.
func LoadConfig() (Config, error) {
	flags, err := getFlags()
	if err != nil {
		return Config{}, fmt.Errorf("config.LoadConfig.getFlags %w", err)
	}

	if checkFlags(flags) {
		return LoadConfigFile(flags["configFile"])
	}

	return LoadConfigEnvironments()
}

// getFlags reads flags.
func getFlags() (map[string]string, error) {
	var flags map[string]string

	rootCmd := new(cobra.Command)
	rootCmd.Use = "myc-device-simulator-app"
	rootCmd.Run = func(cmd *cobra.Command, args []string) {}

	rootCmd.Flags().StringToStringVarP(&flags, "flag", "f", nil, "Flag")

	if err := viper.BindPFlag("flag", rootCmd.Flags().Lookup("flag")); err != nil {
		return flags, fmt.Errorf("config.getFlags.BinPFlag: %w", err)
	}

	if err := rootCmd.Execute(); err != nil {
		return flags, fmt.Errorf("config.getFlags.Execute: %w", err)
	}

	return flags, nil
}

// checkFlags checkout a configFile.
func checkFlags(flags map[string]string) bool {
	if len(flags) == 0 {
		return false
	}

	if flags["configFile"] == "" || len(flags["configFile"]) == 0 {
		return false
	}

	return true
}

// LoadConfigEnvironments loads environment vars into Config.
func LoadConfigEnvironments() (Config, error) {
	var config Config

	viper.SetEnvPrefix("MYC_DEVICES_SIMULATOR")

	// env server.
	_ = viper.BindEnv("HOST")
	_ = viper.BindEnv("HOSTNAME")
	_ = viper.BindEnv("PORT")
	_ = viper.BindEnv("BASEURL")
	_ = viper.BindEnv("SERVERURI")

	// env config sentry.
	_ = viper.BindEnv("SENTRY")
	_ = viper.BindEnv("ENVIRONMENT")
	_ = viper.BindEnv("RELEASE")
	_ = viper.BindEnv("TRACESSAMPLERATE")

	// env DB.
	_ = viper.BindEnv("DBPOSTGRES")
	_ = viper.BindEnv("DBMAXIDLECONNS")
	_ = viper.BindEnv("DBMAXOPENCONNS")
	_ = viper.BindEnv("DBLOGGER")

	// env QUEUE.
	_ = viper.BindEnv("QUEUEHOST")
	_ = viper.BindEnv("QUEUEPORT")
	_ = viper.BindEnv("QUEUECONCURRENCY")

	// env Email.
	_ = viper.BindEnv("POSTMARKTOKEN")
	_ = viper.BindEnv("SMTPHOST")
	_ = viper.BindEnv("SMTPPORT")
	_ = viper.BindEnv("SMTPNETWORK")
	_ = viper.BindEnv("SMTPFROM")
	_ = viper.BindEnv("TEMPLATEFOLDER")

	// env jwt key.
	_ = viper.BindEnv("SECRETKEY")

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("config.LoadConfigEnvironments.Unmarshal: %w", err)
	}

	return config, nil
}

// LoadConfigFile load environment vars into Config which yaml file.
func LoadConfigFile(configFile string) (Config, error) {
	var config Config

	viper.AddConfigPath(".")
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("config.LoadConfigFile.ReadInConfig: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("config.LoadConfigFile.Unmarshal: %w", err)
	}

	return config, nil
}
