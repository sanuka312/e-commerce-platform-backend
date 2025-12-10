package router

import (
	"shophub-backend/auth"

	"github.com/gin-gonic/gin"
)

type CheckoutControllerInterface interface {
	CreateOrder(ctx *gin.Context)
}

func RegisterCheckoutRoutes(router *gin.Engine, controller CheckoutControllerInterface) {
	authMiddleware := auth.AuthMiddleware()
	checkoutGroup := router.Group("/checkout", authMiddleware)
	{
		// Route for placing an order during checkout
		checkoutGroup.POST("/order", controller.CreateOrder)
	}
}
