package controller

import (
	"net/http"
	"shophub-backend/auth"
	"shophub-backend/data"
	"shophub-backend/logger"
	"shophub-backend/model"
	"shophub-backend/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AddressController struct {
	AddressService service.AddressService
}

func NewAddressController(AddressService service.AddressService) *AddressController {
	return &AddressController{
		AddressService: AddressService,
	}
}

func (c *AddressController) GetUserAddresses(ctx *gin.Context) {
	logger.ActInfo("Fetching user addresses")

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
	addresses, err := c.AddressService.GetAddressesByUser(keycloakUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
			Error:            "Internal Server Error",
			ErrorDescription: "Failed to fetch user addresses",
			Details:          err.Error(),
		})
		return
	}
	logger.ActInfo("User addresses fetched successfully")
	ctx.JSON(http.StatusOK, addresses)
}

func (c *AddressController) CreateAddress(ctx *gin.Context) {
	logger.ActInfo("Creating address")

	//extracting keycloak user Id from token claims
	claims := auth.GetClaims(ctx)
	if claims == nil || claims.Sub == "" {
		ctx.JSON(http.StatusUnauthorized, data.ErrorResponse{
			Error:            "unauthorized",
			ErrorDescription: "User not authenticated or missing user ID in token",
		})
		return
	}

	keycloakUserID := claims.Sub
	var req data.CreateAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.ActError("Failed to bind request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Invalid request body. Expected: {line1: string, line2: string, city: string, postal_code: string, country: string}",
			Details:          err.Error(),
		})
		return
	}
	//checking if all field are not null
	if req.Line1 == "" || req.Line2 == "" || req.City == "" || req.PostalCode == "" || req.Country == "" {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "All fields are required",
		})
		return
	}

	//creating address object
	address := &model.Address{
		KeycloakUserID: keycloakUserID,
		Line1:          req.Line1,
		Line2:          req.Line2,
		City:           req.City,
		PostalCode:     req.PostalCode,
		Country:        req.Country,
	}
	//Calling create address service
	if err := c.AddressService.CreateAddress(address); err != nil {
		logger.ActError("Failed to create address", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
			Error:            "Internal Server Error",
			ErrorDescription: "Failed to create address",
			Details:          err.Error(),
		})
		return
	}
	logger.ActInfo("Address created successfully")
	ctx.JSON(http.StatusOK, data.MessageResponse{
		Message: "Address created successfully",
	})
}
