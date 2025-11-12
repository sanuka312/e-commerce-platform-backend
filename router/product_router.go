package router

import (
	"shophub-backend/auth"

	"github.com/gin-gonic/gin"
)

type ProductControllerInterface interface {
	GetAllProducts(ctx *gin.Context)
	GetProductById(ctx *gin.Context)
}

func RegisterProductRoutes(router *gin.Engine, controller ProductControllerInterface) {
	authMiddleware := auth.AuthMiddleware()
	productGroup := router.Group("/products", authMiddleware)
	{

		productGroup.GET("/", controller.GetAllProducts)
		//Route for getting product by ID
		productGroup.GET("/:id", controller.GetProductById)
	}
}
