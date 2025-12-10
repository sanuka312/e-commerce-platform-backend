package router

import (
	"shophub-backend/auth"

	"github.com/gin-gonic/gin"
)

type AddressControllerInterface interface {
	GetUserAddresses(ctx *gin.Context)
	CreateAddress(ctx *gin.Context)
}

func RegisterAddressRoutes(router *gin.Engine, controller AddressControllerInterface) {
	authMiddleware := auth.AuthMiddleware()
	addressGroup := router.Group("/addresses", authMiddleware)
	{
		// Get all addresses for the authenticated user
		addressGroup.GET("/", controller.GetUserAddresses)
		// Create a new address for the authenticated user
		addressGroup.POST("/", controller.CreateAddress)
	}
}
