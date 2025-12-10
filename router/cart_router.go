package router

import (
	"shophub-backend/auth"

	"github.com/gin-gonic/gin"
)

type CartControllerInterface interface {
	GetUserCart(ctx *gin.Context)
	AddItemToCart(ctx *gin.Context)
	ClearCart(ctx *gin.Context)
	RemoveItemFromCart(ctx *gin.Context)
	UpdateCartItemQuantity(ctx *gin.Context)
}

func RegisterCartRoutes(router *gin.Engine, controller CartControllerInterface) {
	authMiddleware := auth.AuthMiddleware()
	cartGroup := router.Group("/cart", authMiddleware)
	{
		cartGroup.GET("/", controller.GetUserCart)
		cartGroup.GET("/:userId", controller.GetUserCart)
		cartGroup.POST("/item", controller.AddItemToCart)
		cartGroup.DELETE("/items", controller.ClearCart)
		cartGroup.DELETE("/item/:itemId", controller.RemoveItemFromCart)
		cartGroup.PATCH("/item/:itemId", controller.UpdateCartItemQuantity)
	}
}
