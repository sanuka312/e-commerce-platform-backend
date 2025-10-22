package logger

import (
	"e-commerce-platform-backend/config"
	"e-commerce-platform-backend/data"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	appLog   *zap.Logger
	actLog   *zap.Logger
	logMutex sync.Mutex
	currDate string
)

var openFile = os.OpenFile
var mkdirAll = os.MkdirAll
var getEnv = os.Getenv
var newTicker = time.NewTicker

func Init() {
	if config.LoadConfig().Env == "test" ||
		config.LoadConfig().Env == "" ||
		os.Getenv("ENV") == "test" ||
		os.Getenv("ENV") == "" {
		appLog = zap.NewNop()
		appLog = zap.NewNop()
		return
	}
	initLogger()
	go monitorDateChange()
}

func initLogger() {
	logMutex.Lock()
	defer logMutex.Unlock()

	fmt.Print("----------------------------------------------------------------------------------->")

	currDate = time.Now().Format(data.DATE_FORMAT_YYYYMMDD)
	appLogPath := getAppLogFilePath()
	actLogPath := getActLogFilePath()

	_ = mkdirAll("logs", os.ModePerm)

	appFile, err := openFile(appLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open app log file: %v", err))
	}

	actFile, err := openFile(actLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open act log file: %v", err))
	}

	atom := zap.NewAtomicLevelAt(getLogLevelFromEnv())

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	appWriter := zapcore.AddSync(appFile)
	actWriter := zapcore.AddSync(actFile)

	var encoder zapcore.Encoder
	if getEnv("ENV") == "production" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		// Add stdout to the loggers
		appWriter = zapcore.NewMultiWriteSyncer(appWriter, zapcore.AddSync(os.Stdout))
		actWriter = zapcore.NewMultiWriteSyncer(actWriter, zapcore.AddSync(os.Stdout))
	}

	appCore := zapcore.NewCore(encoder, appWriter, atom)
	actCore := zapcore.NewCore(encoder, actWriter, atom)

	// Replace existing loggers safely
	if appLog != nil {
		_ = appLog.Sync()
	}
	if actLog != nil {
		_ = actLog.Sync()
	}

	appLog = zap.New(appCore, zap.AddCaller(), zap.AddCallerSkip(1))
	actLog = zap.New(actCore, zap.AddCaller(), zap.AddCallerSkip(1))
}

func getAppLogFilePath() string {
	date := time.Now().Format(data.DATE_FORMAT_YYYYMMDD)
	return filepath.Join("logs", fmt.Sprintf("app-%s.log", date))
}

func getActLogFilePath() string {
	date := time.Now().Format(data.DATE_FORMAT_YYYYMMDD)
	return filepath.Join("logs", fmt.Sprintf("act-%s.log", date))
}

func monitorDateChange() {
	ticker := newTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		today := time.Now().Format(data.DATE_FORMAT_YYYYMMDD)
		if today != currDate {
			logMutex.Lock()
			currDate = today
			logMutex.Unlock()
			initLogger()
		}
	}
}

// Act log functions
func ActDebug(msg string, fields ...zap.Field) {
	logMutex.Lock()
	defer logMutex.Unlock()
	actLog.Debug(msg, fields...)
}

func ActInfo(msg string, fields ...zap.Field) {
	logMutex.Lock()
	defer logMutex.Unlock()
	actLog.Info(msg, fields...)
}

func ActWarn(msg string, fields ...zap.Field) {
	logMutex.Lock()
	defer logMutex.Unlock()
	actLog.Warn(msg, fields...)
}

func ActError(msg string, fields ...zap.Field) {
	logMutex.Lock()
	defer logMutex.Unlock()
	actLog.Error(msg, fields...)
}

// App log functions
func AppDebug(msg string, fields ...zap.Field) {
	logMutex.Lock()
	defer logMutex.Unlock()
	appLog.Debug(msg, fields...)
}
func AppInfo(msg string, fields ...zap.Field) {
	logMutex.Lock()
	defer logMutex.Unlock()
	appLog.Info(msg, fields...)
}

func AppWarn(msg string, fields ...zap.Field) {
	logMutex.Lock()
	defer logMutex.Unlock()
	appLog.Warn(msg, fields...)
}

func AppError(msg string, fields ...zap.Field) {
	logMutex.Lock()
	defer logMutex.Unlock()
	appLog.Error(msg, fields...)
}

func Sync() {
	logMutex.Lock()
	defer logMutex.Unlock()
	_ = actLog.Sync()
	_ = appLog.Sync()
}

func getLogLevelFromEnv() zapcore.Level {
	logLevelStr := getEnv("APP_LOG_LEVEL") // Or any other environment variable name
	switch strings.ToLower(logLevelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel // Default to InfoLevel if not set or invalid
	}
}
