package mlserver

import (
	"fmt"
	"os"
	"strconv"
)

//AgentConfig configuration parameters
type mlConfig struct {
	grpcPort string
}

//update conf instance with default value and environment variables
func (cfg *mlConfig) init(version string) {
	cfg.setDefault()
	cfg.loadConfigUsingEnvVariable()
	cfg.displayConfig(version)
}

//Set default value of configuration
func (cfg *mlConfig) setDefault() {
	cfg.grpcPort = "30107"
}

//Update config with env variables
func (cfg *mlConfig) loadConfigUsingEnvVariable() {
	cfg.grpcPort = cfg.getStringParameter("GRPC_PORT", cfg.grpcPort)
}

//display amp-pilot configuration
func (cfg *mlConfig) displayConfig(version string) {
	fmt.Printf("mLearning version: %v\n", version)
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("Configuration:")
	fmt.Printf("grpc_port: %s\n", cfg.grpcPort)
	fmt.Println("----------------------------------------------------------------------------")
}

//return env variable value, if empty return default value
func (cfg *mlConfig) getStringParameter(envVariableName string, def string) string {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	return value
}

//return env variable value convert to int, if empty return default value
func (cfg *mlConfig) getIntParameter(envVariableName string, def int) int {
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
