package controller

import (
	"net/http"
	"shophub-backend/data"
	"shophub-backend/logger"
	"shophub-backend/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
	userIdParam := ctx.Param("userId")
	if userIdParam == "" {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "bad_request",
			ErrorDescription: "Missing userId path oarameter",
		})
		return
	}

	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "bad_request",
			ErrorDescription: "Invalid user ID",
			Details:          err.Error(),
		})
		return
	}

	cart, err := c.CartService.GetUserCart(uint(userId))
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
	var req data.AddToCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Invalid request body",
			Details:          err.Error(),
		})
		return
	}

	if req.UserID == 0 || req.CartID == 0 || req.ProductID == 0 || req.Quantity <= 0 {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "user_id, cart_id, product_id must be > 0 and quantity must be >= 1",
		})
		return
	}

	if err := c.CartService.AddTOCart(req.UserID, req.CartID, req.ProductID, req.Quantity); err != nil {
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
	idParam := ctx.Param("userId")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Status bad request",
			ErrorDescription: "Missing userId path parameter",
		})
		return
	}

	userId, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Status bad request",
			ErrorDescription: "Invalid userId",
			Details:          err.Error(),
		})
		return
	}

	if err := c.CartService.ClearCart(uint(userId)); err != nil {
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
