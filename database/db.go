package database

import (
	"fmt"
	"log"
	"shophub-backend/config"
	"shophub-backend/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		logger.AppError("Error connecting to the database",
			zap.Error(err),
			zap.String("host", host),
			zap.Int("port", port),
			zap.String("user", user),
			zap.String("database", dbName),
		)
		log.Fatalf("Failed to connect to database.\nError: %v\nConnection details: host=%s port=%d user=%s dbname=%s\n\nPlease check:\n1. PostgreSQL is running\n2. Database '%s' exists\n3. User '%s' has access to the database\n4. Your .env file has correct DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME values",
			err, host, port, user, dbName, dbName, user)
	}

	logger.AppInfo("Connected to PostgreSQL database successfully")
	return db
}
