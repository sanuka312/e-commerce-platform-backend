package database

import (
	"shophub-backend/config"
	"shophub-backend/logger"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"go.uber.org/zap"
)

var openDB = func(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func InitDB() *gorm.DB {

	var (
		host     = config.LoadConfig().DBHost
		port     = config.LoadConfig().DBPort
		user     = config.LoadConfig().DBUser
		password = config.LoadConfig().DBPassword
		dbName   = config.LoadConfig().DBName
		sslMode  = config.LoadConfig().DBSSLMode
	)

	sqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbName,
		sslMode,
	)

	db, err := openDB(sqlInfo)

	if err != nil {
		logger.AppError("Error connecting to the database", zap.Error(err))
		log.Fatal("Failed to connect to database. Please check your database configuration and ensure PostgreSQL is running.")
	}

	logger.AppInfo("Connected to PostgreSQL database successfully")
	return db
}
