package database

import (
	"e-commerce-platform-backend/config"
	"e-commerce-platform-backend/logger"
	"fmt"
	"log"

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
		log.Fatal("Error connecting the database")
	}

	logger.AppInfo("connected to PostgreSQL database successfully")
	return db
}
