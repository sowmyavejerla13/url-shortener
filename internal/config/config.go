package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppName    string
	AppEnv     string
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
	JWTExpiry  string
}

func LoadConfig() (*Config, error) {
	// Load .env file (ignore if it doesn't exist)
	_ = godotenv.Load()

	// Read environment variables
	viper.AutomaticEnv()

	cfg := &Config{
		AppName:    viper.GetString("APP_NAME"),
		AppEnv:     viper.GetString("APP_ENV"),
		AppPort:    viper.GetString("APP_PORT"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBSSLMode:  viper.GetString("DB_SSL_MODE"),
		JWTSecret:  viper.GetString("JWT_SECRET"),
		JWTExpiry:  viper.GetString("JWT_EXPIRY"),
	}

	required := map[string]string{
		"APP_NAME":    cfg.AppName,
		"APP_ENV":     cfg.AppEnv,
		"APP_PORT":    cfg.AppPort,
		"DB_HOST":     cfg.DBHost,
		"DB_PORT":     cfg.DBPort,
		"DB_USER":     cfg.DBUser,
		"DB_PASSWORD": cfg.DBPassword,
		"DB_NAME":     cfg.DBName,
		"DB_SSL_MODE": cfg.DBSSLMode,
		"JWT_SECRET":  cfg.JWTSecret,
		"JWT_EXPIRY":  cfg.JWTExpiry,
	}

	for key, value := range required {
		if value == "" {
			return nil, fmt.Errorf("%s is required", key)
		}
	}

	return cfg, nil
}
