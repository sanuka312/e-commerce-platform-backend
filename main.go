package main

import (
	"net/http"
	"os"
	"shophub-backend/config"
	"shophub-backend/controller"
	"shophub-backend/database"
	"shophub-backend/logger"
	"shophub-backend/migration"
	"shophub-backend/repository"
	"shophub-backend/router"
	"shophub-backend/service"
	"shophub-backend/utils"
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

	logger.Init()

	logger.AppInfo("Starting the backend service")

	logger.AppInfo("Loading configuration...")
	if os.Getenv("ENV") != "production" {
		config.LoadEnv()
	}

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

	//Initializing the repository files
	cartRepository := repository.NewCartRepository(pgDb)
	productRepository := repository.NewProductRepository(pgDb)
	orderRepository := repository.NewOrderRepository(pgDb)
	paymentRepository := repository.NewPaymentRepositoryImpl(pgDb)

	cartService, err := service.NewCartServiceImpl(cartRepository, productRepository)
	if err != nil {
		logger.ActError("Failed to initialize the cart service", zap.Error(err))
		return
	}

	productService, err := service.NewProductServiceImpl(productRepository)
	if err != nil {
		logger.ActError("Failed to initialize the product service", zap.Error(err))
		return
	}

	orderService, err := service.NewOrderServiceImpl(orderRepository, productRepository, cartRepository, paymentRepository)
	if err != nil {
		logger.ActError("Failed to initialize the order service", zap.Error(err))
		return
	}

	paymentService, err := service.NewPaymentServiceImpl(paymentRepository, orderRepository)
	if err != nil {
		logger.ActError("Failed to initialize the payment service", zap.Error(err))
		return
	}

	//Initializing the controllers
	cartController := controller.NewCartController(cartService)
	productController := controller.NewProductController(productService)
	orderController := controller.NewOrderController(orderService)
	paymentController := controller.NewPaymentController(paymentService)

	//Create gin router
	r := gin.Default()

	// Serve static files (images) from assets directory
	r.Static("/assets", "./assets")

	//Register routes
	router.RegisterCartRoutes(r, cartController)
	router.RegisterProductRoutes(r, productController)
	router.RegisterOrderRoutes(r, orderController)
	router.RegisterPaymentRoutes(r, paymentController)

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
	err = server.ListenAndServe()
	utils.ErrorPanic(err)
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
