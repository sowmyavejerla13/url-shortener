package config

import (
	"log"

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

func LoadConfig() *Config {
	_ = godotenv.Load()
	viper.AutomaticEnv()

	cfg := &Config{
		AppName:    viper.GetString("APP_NAME"),
		AppPort:    viper.GetString("APP_PORT"),
		AppEnv:     viper.GetString("APP_ENV"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBSSLMode:  viper.GetString("SSL_MODE"),
		JWTSecret:  viper.GetString("JWT_SECRET"),
		JWTExpiry:  viper.GetString("JWT_EXPIRY"),
	}

	log.Println("Configuration Load")

	return cfg

}
