package controller

import (
	"e-commerce-platform-backend/data"
	"e-commerce-platform-backend/logger"
	"e-commerce-platform-backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	ProductService service.ProductService
}

func NewProductController(ProductService service.ProductService) *ProductController {
	return &ProductController{
		ProductService: ProductService,
	}
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	logger.ActInfo("Fetching all products")
	products, err := c.ProductService.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, data.ErrorResponse{
			Error:            "Internal Server Error",
			ErrorDescription: "Failed to fetch the products",
			Details:          err.Error(),
		})
		return
	}
	logger.ActInfo("Products fetched successfully")
	ctx.JSON(http.StatusOK, products)
}

func (c *ProductController) GetProductById(ctx *gin.Context) {
	logger.ActInfo("Fetching product by ID")
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Invalid product ID ",
			Details:          err.Error(),
		})
		return
	}
	product, err := c.ProductService.GetProductById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, data.ErrorResponse{
			Error:            "Bad Request",
			ErrorDescription: "Product not found",
			Details:          err.Error(),
		})
		return
	}
	logger.ActInfo("Product fetched succcessfully")
	ctx.JSON(http.StatusOK, product)
}
