package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	SecretKey string         `mapstructure:"SECRET_KEY"`
	DBConfig  DatabaseConfig `mapstructure:"database"`
	Server    ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	BaseURL string `mapstructure:"base_url"`
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("server.port", "SERVER_PORT")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file (config.yaml) not found, relying on environment variables.")
		} else {
			log.Fatalf("Failed to read config file: %v", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "debug"
	}
	if cfg.DBConfig.SSLMode == "" {
		cfg.DBConfig.SSLMode = "disable"
	}
	if cfg.Server.BaseURL == "" {
		cfg.Server.BaseURL = "http://localhost:8080"
	}

	if cfg.SecretKey == "" {
		log.Fatal("SECRET_KEY is not set in config.")
	}
	if cfg.Server.Port == "" {
		log.Fatal("SERVER_PORT is not set in config.")
	}

	return &cfg
}
