package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found")
	}
}

func LoadConfig() *Config {

	return &Config{
		DBHost:     Getenv("ENV", "localhost"),
		DBPort:     GetenvAsInt("DB_PORT", 5432),
		DBUser:     Getenv("DB_USER", ""),
		DBPassword: Getenv("DB_PASSWORD", ""),
		DBName:     Getenv("DB_NAME", "shopping"),
		DBSSLMode:  Getenv("DB_SSLMODE", "disable"),
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
