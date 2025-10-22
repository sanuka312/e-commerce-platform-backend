package main

import (
	"e-commerce-platform-backend/config"
	"e-commerce-platform-backend/database"
	"e-commerce-platform-backend/logger"
	"e-commerce-platform-backend/migration"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

var allowedOrigins []string

func main() {

	defer logger.Sync()

	if os.Getenv("ENV") != "production" {
		config.LoadEnv()
	}

	logger.Init()
	defer logger.Sync()

	logger.AppInfo("Starting the backend service")

	logger.AppInfo("Loading configuration...")

	//Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	//Load Origins
	processAllowedOrigins()

	logger.AppInfo("Loading database configurations")
	pgDb := database.InitDB()

	if err := migration.Migrate(pgDb); err != nil {
		logger.AppError("Migration failed", zap.Error(err))
	}

	r := gin.Default()

	// Enable CORS for all origins
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-User-ID"},
		AllowCredentials: true,
	}).Handler(r)

	server := &http.Server{
		Addr:           ":" + strconv.Itoa(config.LoadConfig().Port),
		Handler:        corsHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.AppInfo("Server started on port " + strconv.Itoa(config.LoadConfig().Port))
	server.ListenAndServe()
}

func processAllowedOrigins() {
	origins := config.LoadConfig().AllowedOrigins
	if origins == "" {
		logger.AppError("Allowed origins not set")

	}
	allowedOrigins = strings.Split(origins, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}
}
