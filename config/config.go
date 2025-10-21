package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	AppLogLevel string
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
}

func LoadEnv() {
	var err error
	if os.Getenv("ENV") == "test" {
		envFile := filepath.Join("..", ".env.test")
		err = godotenv.Load(envFile)
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func LoadConfig() *Config {

	return &Config{
		AppLogLevel: Getenv("APP_LOG_LEVEL", "INFO"),
		Env:         Getenv("ENV", "development"),
		DBHost:      Getenv("DB_HOST", "localhost"),
		DBPort:      GetenvAsInt("DB_PORT", 5432),
		DBUser:      Getenv("DB_USER", "postgres"),
		DBPassword:  Getenv("DB_PASSWORD", ""),
		DBName:      Getenv("DB_NAME", "shopping_website"),
		DBSSLMode:   Getenv("DB_SSLMODE", "disable"),
	}

}

func Getenv(key, fallBack string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallBack
}

func GetenvAsInt(key string, fallBack int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallBack
}
