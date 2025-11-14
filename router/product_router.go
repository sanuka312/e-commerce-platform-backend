package router

import (
	"github.com/gin-gonic/gin"
)

type ProductControllerInterface interface {
	GetAllProducts(ctx *gin.Context)
	GetProductById(ctx *gin.Context)
	GetProductBySlug(ctx *gin.Context)
}

func RegisterProductRoutes(router *gin.Engine, controller ProductControllerInterface) {

	productGroup := router.Group("/products")
	{

		productGroup.GET("/", controller.GetAllProducts)
		//Route for getting product by ID
		productGroup.GET("/:id", controller.GetProductById)
		//Route foe getting product by product slug
		productGroup.GET("/slug/:productSlug", controller.GetProductBySlug)
	}
}
