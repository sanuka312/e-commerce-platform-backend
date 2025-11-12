package migration

import (
	"shophub-backend/logger"
	"shophub-backend/model"
)

type DBMigrator interface {
	AutoMigrate(dst ...interface{}) error
}

func Migrate(db DBMigrator) error {
	logger.AppInfo("Database Migration")
	return db.AutoMigrate(
		&model.User{},
	)
}
