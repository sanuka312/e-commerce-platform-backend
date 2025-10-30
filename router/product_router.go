package router

import "github.com/gin-gonic/gin"

type ProductControllerInterface interface {
	GetAllProducts(ctx *gin.Context)
	GetProductById(ctx *gin.Context)
}

func RegisterProductRoutes(router *gin.Engine, controller ProductControllerInterface) {
	productGroup := router.Group("/products")
	{
		productGroup.GET("/", controller.GetAllProducts)
		productGroup.GET("/:id", controller.GetProductById)
	}
}
