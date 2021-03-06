package config

import (
	"log"

	"github.com/spf13/viper"
)

var Config *Configuration

// Struct of Configuration instance.
// It include Database and Server configuration
type Configuration struct {
	Server        ServerConnection
	Database      DatabaseConfiguration
	Database_Test DatabaseTestConfiguration
}

// Struct of Server Configuration instance.
type ServerConnection struct {
	Port             string `mapstructure:"SERVER_PORT"`
	Secret           string `mapstructure:"SERVER_SECRET"`
	Endpoint         string `mapstructure:"SERVER_ENDPOINT"`
	Mode             string `mapstructure:"SERVER_MODE"`
	Name             string `mapstructure:"SERVER_NAME"`
	ExpiresHour      int    `mapstructure:"SERVER_EXPIRES_HOUR"`
	UrlCharacterLong int    `mapstructure:"URL_SHORTENER_CHARACTER_LONG"`
}

// Struct of Database Configuration instance.
type DatabaseConfiguration struct {
	Driver       string `mapstructure:"DATABASE_DRIVER"`
	Dbname       string `mapstructure:"DATABASE_NAME"`
	Username     string `mapstructure:"DATABASE_USERNAME"`
	Password     string `mapstructure:"DATABASE_PASSWORD"`
	Host         string `mapstructure:"DATABASE_HOST"`
	Port         string `mapstructure:"DATABASE_PORT"`
	MaxLifetime  int    `mapstructure:"DATABASE_MAX_LIFETIME"`
	MaxOpenConns int    `mapstructure:"DATABASE_MAX_OPEN_CONNS"`
	MaxIdleConns int    `mapstructure:"DATABASE_MAX_IDLE_CONNS"`
	SslMode      string `mapstructure:"DATABASE_SSL_MODE"`
}

// Struct of Database for Testing Configuration instance
type DatabaseTestConfiguration struct {
	Driver       string `mapstructure:"DATABASE_TEST_DRIVER"`
	Dbname       string `mapstructure:"DATABASE_TEST_NAME"`
	Username     string `mapstructure:"DATABASE_TEST_USERNAME"`
	Password     string `mapstructure:"DATABASE_TEST_PASSWORD"`
	Host         string `mapstructure:"DATABASE_TEST_HOST"`
	Port         string `mapstructure:"DATABASE_TEST_PORT"`
	MaxLifetime  int    `mapstructure:"DATABASE_TEST_MAX_LIFETIME"`
	MaxOpenConns int    `mapstructure:"DATABASE_TEST_MAX_OPEN_CONNS"`
	MaxIdleConns int    `mapstructure:"DATABASE_TEST_MAX_IDLE_CONNS"`
	SslMode      string `mapstructure:"DATABASE_TEST_SSL_MODE"`
}

// Setup the configuration
func Setup(configPath string) {
	var (
		serverConfiguration       ServerConnection
		databaseConfiguration     DatabaseConfiguration
		databaseTestConfiguration DatabaseTestConfiguration
	)

	viper.SetConfigFile(configPath)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	unmarshalConfiguration(&databaseConfiguration)
	unmarshalConfiguration(&databaseTestConfiguration)
	unmarshalConfiguration(&serverConfiguration)

	configuration := Configuration{
		Database:      databaseConfiguration,
		Database_Test: databaseTestConfiguration,
		Server:        serverConfiguration,
	}

	Config = &configuration
}

// Helper to unmarshal
func unmarshalConfiguration(configuration interface{}) {
	err := viper.Unmarshal(configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}

// GetConfig return the configuration instance
func GetConfig() *Configuration {
	return Config
}
