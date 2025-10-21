package main

import (
	"e-commerce-platform-backend/config"
	"e-commerce-platform-backend/database"
	"e-commerce-platform-backend/logger"
	"e-commerce-platform-backend/migration"
	"os"

	"go.uber.org/zap"
)

func main() {

	defer logger.Sync()

	logger.AppInfo("Starting the backend service")

	logger.AppInfo("Loading configuration...")
	if os.Getenv("ENV") != "production" {
		config.LoadEnv()
	}

	logger.AppInfo("Loading database configurations")
	pgDb := database.InitDB()

	if err := migration.Migrate(pgDb); err != nil {
		logger.AppError("Migration failed", zap.Error(err))
	}

}
