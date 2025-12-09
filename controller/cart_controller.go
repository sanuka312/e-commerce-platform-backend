package controller

import (
	"net/http"
	"shophub-backend/auth"
	"shophub-backend/data"
	"shophub-backend/logger"
	"shophub-backend/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CartController struct {
	CartService service.CartService
}

func NewCartController(CartService service.CartService) *CartController {
	return &CartController{
		CartService: CartService,
	}
}

func (c *CartController) GetUserCart(ctx *gin.Context) {
	logger.ActInfo("Fetching User Cart")

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
	cart, err := c.CartService.GetUserCart(keycloakUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
			Error:            "Internal Server Error",
			ErrorDescription: "Failed to fetch the user cart",
			Details:          err.Error(),
		})
		return
	}
	logger.ActInfo("User cart fetched successfully")
	ctx.JSON(http.StatusOK, cart)
}

func (c *CartController) AddItemToCart(ctx *gin.Context) {
	logger.ActInfo("Adding Items to the cart")

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

	var req data.AddToCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.ActError("Failed to bind request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Invalid request body. Expected: {product_id: number, quantity: number}",
			Details:          err.Error(),
		})
		return
	}

	if req.ProductID == 0 || req.Quantity <= 0 {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "product_id must be > 0 and quantity must be >= 1",
		})
		return
	}

	if err := c.CartService.AddTOCart(keycloakUserID, req.ProductID, req.Quantity); err != nil {
		if strings.Contains(err.Error(), "insufficient stock") {
			ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
				Error:            "Bad Request",
				ErrorDescription: err.Error(),
			})
		} else if strings.Contains(err.Error(), "product not found") {
			ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
				Error:            "Not Found",
				ErrorDescription: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
				Error:            "Internal Server Error",
				ErrorDescription: "Failed add item to the user cart",
				Details:          err.Error(),
			})
		}
		return
	}
	logger.ActInfo("Item added to the cart successfully")
	ctx.JSON(http.StatusOK, data.MessageResponse{
		Message: "Item added to the cart successfully",
	})
}

func (c *CartController) ClearCart(ctx *gin.Context) {
	logger.ActInfo("Clearing the user cart")

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

	if err := c.CartService.ClearCart(keycloakUserID); err != nil {
		ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
			Error:            "Internal Server Error",
			ErrorDescription: "Failed to clear user cart",
			Details:          err.Error(),
		})
		return
	}
	logger.ActInfo("User cart cleared successfully")
	ctx.JSON(http.StatusOK, data.MessageResponse{
		Message: "User cart cleared",
	})
}

func (c *CartController) RemoveItemFromCart(ctx *gin.Context) {
	logger.ActInfo("Removing item from user cart")

	itemIdParam := ctx.Param("itemId")
	if itemIdParam == "" {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Status bad request",
			ErrorDescription: "Missing itemId path parameter",
		})
		return
	}

	itemId, err := strconv.ParseUint(itemIdParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Status bad request",
			ErrorDescription: "Invalid itemId",
			Details:          err.Error(),
		})
		return
	}

	if err := c.CartService.RemoveItemFromCart(uint(itemId)); err != nil {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Failed to remove the item from the cart",
			Details:          err.Error(),
		})
		return
	}

	logger.ActInfo("Item removed from cart successfully")
	ctx.JSON(http.StatusOK, data.MessageResponse{
		Message: "Item removed from cart successfully",
	})
}

func (c *CartController) UpdateCartItemQuantity(ctx *gin.Context) {
	logger.ActInfo("Updating cart item quantity")
	itemIdParam := ctx.Param("itemId")
	if itemIdParam == "" {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Missing itemId path parameter",
		})
		return
	}

	itemId, err := strconv.ParseUint(itemIdParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Invalid itemId",
			Details:          err.Error(),
		})
		return
	}

	var req data.UpdateCartItemQuantityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Invalid request body",
			Details:          err.Error(),
		})
		return
	}

	if req.Quantity <= 0 {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Quantity must be greater than zero",
		})
		return
	}

	if err := c.CartService.UpdateCartItemQuantity(uint(itemId), req.Quantity); err != nil {
		ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
			Error:            "Internal Server Error",
			ErrorDescription: "Failed to update cart item quantity",
			Details:          err.Error(),
		})
		return
	}
	logger.ActInfo("Cart item quantity updated successfully")
	ctx.JSON(http.StatusOK, data.MessageResponse{
		Message: "Cart item quantity updated successfully",
	})
}
