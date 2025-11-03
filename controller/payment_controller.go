package controller

import (
	"e-commerce-platform-backend/data"
	"e-commerce-platform-backend/logger"
	"e-commerce-platform-backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	PaymentService service.PaymentService
}

func NewPaymentController(PaymentService service.PaymentService) *PaymentController {
	return &PaymentController{
		PaymentService: PaymentService,
	}
}

func (c *PaymentController) GetPaymentByOrderId(ctx *gin.Context) {
	logger.ActInfo("Fetching payment Order ID's")
	idParam := ctx.Param("orderId")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Missing OrderID parameter",
		})
		return
	}

	orderId, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
			Error:            "Internal Server Error",
			ErrorDescription: "Invalid OrderID",
			Details:          err.Error(),
		})
		return
	}

	payment, err := c.PaymentService.GetPaymentByOrderId(uint(orderId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
			Error:            "Internal Server Error",
			ErrorDescription: "Failed to fetch payment",
			Details:          err.Error(),
		})
		return
	}

	logger.ActInfo("Payments fetched successfully")
	ctx.JSON(http.StatusOK, payment)
}
