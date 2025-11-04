package router

import "github.com/gin-gonic/gin"

type PaymentControlInterface interface {
	GetPaymentByOrderId(ctx *gin.Context)
	ProcessPayment(ctx *gin.Context)
}

func RegisterPaymentRoutes(router *gin.Engine, controller PaymentControlInterface) {
	paymentGroup := router.Group("/payments")
	{
		//get payment details for an order
		paymentGroup.GET("/order/:orderId", controller.GetPaymentByOrderId)

		//Processing the payment for an order
		paymentGroup.POST("/order/:orderId/process")
	}
}
