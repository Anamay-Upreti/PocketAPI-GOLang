package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	AppName     string
	Port        string
	Environment string
	JWTSecret   string
	DatabaseURL string
}

var AppConfig Config

func LoadConfig() {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	AppConfig = Config{
		AppName:     getEnv("APP_NAME", "PocketAPI"),
		Port:        getEnv("PORT", "8000"),
		Environment: getEnv("ENV", "development"),
		JWTSecret:   getEnv("JWT_SECRET", "secret"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key string, defaultValue string) string {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}