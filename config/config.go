package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env  string
	Port int

	AllowedOrigins  string
	DBHost          string
	DBPort          int
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	IdpBaseUrl      string
	IdpRealm        string
	IdpClientSecret string
	IdpClientId     string
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
		log.Printf("Warning: Could not load .env file: %v. Using default values or environment variables.", err)
	}
}

func LoadConfig() *Config {

	return &Config{
		Env:  Getenv("ENV", "development"),
		Port: GetenvAsInt("PORT", 9002),

		AllowedOrigins:  Getenv("ALLOWED_ORIGINS", ""),
		DBHost:          Getenv("DB_HOST", "localhost"),
		DBPort:          GetenvAsInt("DB_PORT", 5432),
		DBUser:          Getenv("DB_USER", "postgres"),
		DBPassword:      Getenv("DB_PASSWORD", ""),
		DBName:          Getenv("DB_NAME", "shophub_website"),
		DBSSLMode:       Getenv("DB_SSLMODE", "disable"),
		IdpBaseUrl:      Getenv("IDP_BASE_URL", ""),
		IdpRealm:        Getenv("IDP_REALM", ""),
		IdpClientId:     Getenv("IDP_CLIENT_ID", ""),
		IdpClientSecret: Getenv("IDP_CLIENT_SECRET", ""),
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
