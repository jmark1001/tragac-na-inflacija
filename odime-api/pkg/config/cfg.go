package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config holds the application configuration values.
type Config struct {
	AppName          string `mapstructure:"app_name"`
	ServerPort       int    `mapstructure:"server_port"`
	MQHost           string `mapstructure:"mq_host"`
	MQPort           int    `mapstructure:"mq_port"`
	MQUser           string `mapstructure:"mq_user"`
	MQPassword       string `mapstructure:"mq_password"`
	MQPendingQueue   string `mapstructure:"mq_pending_queue"`
	MQProcessedQueue string `mapstructure:"mq_processed_queue"`
	KafkaBroker      string `mapstructure:"kafka_broker"`
	KafkaTopic       string `mapstructure:"kafka_topic"`
	DbHost           string `mapstructure:"db_path"`
	DbPort           int    `mapstructure:"db_port"`
	DbUser           string `mapstructure:"db_user"`
	DbPassword       string `mapstructure:"db_password"`
	DbName           string `mapstructure:"db_name"`
}

// LoadConfig reads the configuration from a file or environment variables.
func LoadConfig() (*Config, error) {
	// Set default values for environment variables
	viper.SetDefault("app_name", "MyApp")
	viper.SetDefault("server_port", 8080)

	// Tell Viper to read environment variables and config file
	viper.AutomaticEnv()

	// Optionally, you can specify the config file path or format
	viper.SetConfigFile(".env") // Or use a specific config file path

	// Read the config file (if any)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config file found, using environment variables. Error: %v", err)
	}

	// Unmarshal the configuration values into the Config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %v", err)
	}

	// Return the config struct
	return &cfg, nil
}
