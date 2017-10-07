package main

import (
	"fmt"
	"os"
	"strconv"
)

//CliConfig configuration parameters
type CliConfig struct {
	serverAddress string
	colorTheme    string
}

//update conf instance with default value and environment variables
func (cfg *CliConfig) init(version string, build string) {
	cfg.setDefault()
	cfg.loadConfigUsingEnvVariable()
	//cfg.displayConfig(version, build)
}

//Set default value of configuration
func (cfg *CliConfig) setDefault() {
	cfg.serverAddress = "127.0.0.1:30107"
	cfg.colorTheme = "dark"
}

//Update config with env variables
func (cfg *CliConfig) loadConfigUsingEnvVariable() {
	cfg.serverAddress = cfg.getStringParameter("SERVER_ADDRESS", cfg.serverAddress)
	cfg.colorTheme = cfg.getStringParameter("COLOR_THEME", cfg.colorTheme)
}

//display amp-pilot configuration
func (cfg *CliConfig) displayConfig(version string, build string) {
	fmt.Printf("mlearning version: %v build: %s\n", version, build)
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("Configuration:")
	fmt.Printf("mlearning address: %s\n", cfg.serverAddress)
}

//return env variable value, if empty return default value
func (cfg *CliConfig) getStringParameter(envVariableName string, def string) string {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	return value
}

//return env variable value convert to int, if empty return default value
func (cfg *CliConfig) getIntParameter(envVariableName string, def int) int {
	value := os.Getenv(envVariableName)
	if value != "" {
		ivalue, err := strconv.Atoi(value)
		if err != nil {
			return def
		}
		return ivalue
	}
	return def
}
