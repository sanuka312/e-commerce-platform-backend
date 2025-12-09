package controller

import (
	"net/http"
	"shophub-backend/auth"
	"shophub-backend/data"
	"shophub-backend/logger"
	"shophub-backend/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CheckoutController struct {
	CheckoutService service.CheckoutService
}

func NewCheckoutController(CheckoutService service.CheckoutService) *CheckoutController {
	return &CheckoutController{
		CheckoutService: CheckoutService,
	}
}

func (c *CheckoutController) CreateOrder(ctx *gin.Context) {
	logger.ActInfo("Placing order through checkout")

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

	// Bind request body
	var req data.PlaceOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.ActError("Failed to bind request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Invalid request body. Expected: {payment_method: string, address: {line1: string, line2: string, city: string, postal_code: string, country: string}}",
			Details:          err.Error(),
		})
		return
	}

	// Validate payment method
	if req.PaymentMethod == "" {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Payment method is required",
		})
		return
	}

	// Validate address fields
	if req.Address.Line1 == "" || req.Address.Line2 == "" || req.Address.City == "" || req.Address.PostalCode == "" || req.Address.Country == "" {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "All address fields are required",
		})
		return
	}

	// Call checkout service to place order
	order, err := c.CheckoutService.PlaceOrder(keycloakUserID, req.PaymentMethod, req.Address)
	if err != nil {
		logger.ActError("Failed to place order", zap.Error(err))
		if err.Error() == "cart not found" || err.Error() == "cart is empty" {
			ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
				Error:            "Bad Request",
				ErrorDescription: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
				Error:            "Internal Server Error",
				ErrorDescription: "Failed to place order",
				Details:          err.Error(),
			})
		}
		return
	}

	logger.ActInfo("Order placed successfully through checkout")
	ctx.JSON(http.StatusOK, order)
}
