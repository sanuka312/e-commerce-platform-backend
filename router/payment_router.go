package router

import (
	"shophub-backend/auth"

	"github.com/gin-gonic/gin"
)

type PaymentControlInterface interface {
	GetPaymentByOrderId(ctx *gin.Context)
	ProcessPayment(ctx *gin.Context)
}

func RegisterPaymentRoutes(router *gin.Engine, controller PaymentControlInterface) {
	authMiddleware := auth.AuthMiddleware()
	paymentGroup := router.Group("/payments", authMiddleware)
	{
		//get payment details for an order
		paymentGroup.GET("/order/:orderId", controller.GetPaymentByOrderId)

		//Processing the payment for an order
		paymentGroup.POST("/order/:orderId/process", controller.ProcessPayment)
	}
}
