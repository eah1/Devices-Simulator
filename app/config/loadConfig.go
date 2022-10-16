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
		fmt.Println(err)
	}

	if checkFlags(flags) {
		return LoadConfigFile(flags["configFile"])
	}

	return LoadConfigEnvironments()
}

// getFlags reads flags.
func getFlags() (map[string]string, error) {
	var flags map[string]string

	//nolint: exhaustivestruct
	rootCmd := &cobra.Command{
		Use: "myc-cloud-app",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	rootCmd.Flags().StringToStringVarP(&flags, "flag", "f", nil, "Flag")

	if err := viper.BindPFlag("flag", rootCmd.Flags().Lookup("flag")); err != nil {
		fmt.Printf("error in BindPFlag, %v\n", err)

		return flags, err
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error in executet, %v\n", err)

		return flags, err
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

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("unable to decode into struct, %v\n", err)

		return config, err
	}

	return config, nil
}

// LoadConfigFile load environment vars into Config which yaml file.
func LoadConfigFile(configFile string) (Config, error) {
	var config Config

	viper.AddConfigPath(".")
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("this file not exist, %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("unable to decode into struct, %v\n", err)

		return config, err
	}

	return config, nil
}
