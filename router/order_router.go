package router

import (
	"shophub-backend/auth"

	"github.com/gin-gonic/gin"
)

type OrderControllerInterface interface {
	CreateOrder(ctx *gin.Context)
	GetOrderByUser(ctx *gin.Context)
}

// registering order route nested with payment route
func RegisterOrderRoutes(router *gin.Engine, controller OrderControllerInterface) {
	authMiddleware := auth.AuthMiddleware()
	orderGroup := router.Group("/orders", authMiddleware)
	{
		//Creates a new order for a user
		orderGroup.POST("/", controller.CreateOrder)

		//Get all order for a user
		orderGroup.GET("/user", controller.GetOrderByUser)
	}

}
