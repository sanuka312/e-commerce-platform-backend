package utils

import (
	"shophub-backend/logger"

	"go.uber.org/zap"
)

var AppErrorFunc = logger.AppError

var ErrorPanic = func(err error) {
	if err != nil {
		AppErrorFunc("An error occurred", zap.Error(err))
		panic(err)
	}
}
