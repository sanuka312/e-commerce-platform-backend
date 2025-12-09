package controller

import (
	"net/http"
	"shophub-backend/auth"
	"shophub-backend/data"
	"shophub-backend/logger"
	"shophub-backend/service"

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
	
	// Extract Keycloak user ID from token claims
	claims := auth.GetClaims(ctx)
	if claims == nil || claims.Sub == "" {
		ctx.JSON(http.StatusUnauthorized, data.ErrorResponse{
			Error:            "unauthorized",
			ErrorDescription: "User not authenticated or missing user ID in token",
		})
		return
	}

	keycloakUserID := claims.Sub
	
	//calling the create order service
	order, err := c.OrderService.CreateOrder(keycloakUserID)
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
	
	// Extract Keycloak user ID from token claims
	claims := auth.GetClaims(ctx)
	if claims == nil || claims.Sub == "" {
		ctx.JSON(http.StatusUnauthorized, data.ErrorResponse{
			Error:            "unauthorized",
			ErrorDescription: "User not authenticated or missing user ID in token",
		})
		return
	}

	keycloakUserID := claims.Sub

	//calling the order service
	orders, err := c.OrderService.GetOrderByUser(keycloakUserID)
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
