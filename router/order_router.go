package router

import "github.com/gin-gonic/gin"

type OrderControllerInterface interface {
	CreateOrder(ctx *gin.Context)
	GetOrderByUser(ctx *gin.Context)
}

// registering order route nested with payment route
func RegisterOrderRoutes(router *gin.Engine, controller OrderControllerInterface) {
	orderGroup := router.Group("/orders")
	{
		//Creates a new order for a user
		orderGroup.POST("/:user_id", controller.CreateOrder)

		//Get all order for a user
		orderGroup.GET("/user/:user_id", controller.GetOrderByUser)
	}

}
