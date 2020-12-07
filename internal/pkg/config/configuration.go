package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config stores the server's configuration
var Config *Configuration

// Configuration struct captures "ServerConfiguration"
type Configuration struct {
	Server ServerConfiguration
}

// ServerConfiguration struct captures "Port" and "Mode"
type ServerConfiguration struct {
	Port string
	Mode string
}

// Setup initialize configuration
func Setup(configPath string) {
	var configuration *Configuration

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	Config = configuration
}

// GetConfig helps you to get configuration data
func GetConfig() *Configuration {
	return Config
}
