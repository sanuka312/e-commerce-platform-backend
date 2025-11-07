package router

import (
	"e-commerce-platform-backend/auth"

	"github.com/gin-gonic/gin"
)

type CartControllerInterface interface {
	GetUserCart(ctx *gin.Context)
	AddItemToCart(ctx *gin.Context)
	ClearCart(ctx *gin.Context)
}

func RegisterCartRoutes(router *gin.Engine, controller CartControllerInterface) {
	authMiddleware := auth.AuthMiddleware()
	cartGroup := router.Group("/cart", authMiddleware)
	{
		cartGroup.GET("/:userId", controller.GetUserCart)
		cartGroup.POST("/item", controller.AddItemToCart)
		cartGroup.DELETE("/:userId/items", controller.ClearCart)
	}
}
