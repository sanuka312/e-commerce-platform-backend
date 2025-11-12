package controller

import (
	"shophub-backend/data"
	"shophub-backend/logger"
	"shophub-backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	OrderService service.OrderService
}

func NewOrderController(OrderService service.OrderService) *OrderController {
	return &OrderController{
		OrderService: OrderService,
	}
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	logger.ActInfo("Creating order")
	//Getting user id as a parameter
	userIdParam := ctx.Param("user_id")
	//handling the error if the user id is not provided
	if userIdParam == "" {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Missing UserID parameter",
		})
		return
	}
	//converting the user id parameter to uint
	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	//if the conversion fails, return a bad request error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Invalid UserId parameter",
			Details:          err.Error(),
		})
		return
	}
	//calling the create order service
	order, err := c.OrderService.CreateOrder(uint(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
			Error:            "Internal Server Error",
			ErrorDescription: "Failed to create order",
			Details:          err.Error(),
		})
		return
	}
	//returning the order
	logger.ActInfo("Order created successfully")
	ctx.JSON(http.StatusOK, order)
}

func (c *OrderController) GetOrderByUser(ctx *gin.Context) {
	logger.ActInfo("Fetching orders by user")
	//getting user id as parameter
	userIdParam := ctx.Param("user_id")
	//handling the error if the userId parameter is missing
	if userIdParam == "" {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Missing UserID parameter",
		})
		return
	}

	//converting userId param to uint
	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	//handling th error if the conversion fails
	if err != nil {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Invalid UserId parameter",
			Details:          err.Error(),
		})
		return
	}

	//calling the order service
	orders, err := c.OrderService.GetOrderByUser(uint(userId))
	//handling the error if order service fails
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
			Error:            "Internal Server Error",
			ErrorDescription: "Failed to fetch orders",
			Details:          err.Error(),
		})
		return
	}
	//returning the orders
	logger.ActInfo("Orders fetched successfully")
	ctx.JSON(http.StatusOK, orders)
}
