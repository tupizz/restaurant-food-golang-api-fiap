package handler

import (
	"net/http"
	"strconv"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/mappers"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"

	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	GetProducts(c *gin.Context)
}

type productHandler struct {
	getProductsUseCase usecase.GetProductsUseCase
}

func NewProductHandler(getProductsUseCase usecase.GetProductsUseCase) ProductHandler {
	return &productHandler{getProductsUseCase: getProductsUseCase}
}

// GetProducts godoc
// @Summary      Get Products
// @Description  Get Products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page     query     int  false  "Page number"
// @Param        pageSize query     int  false  "Page size"
// @Param        category query     string  false  "Category"
// @Success      200      {array}  dto.ProductOutput
// @Failure      400      {object}  handler.ErrorResponse
// @Failure      500      {object}  handler.ErrorResponse
// @Router       /products [get]
func (h *productHandler) GetProducts(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
	}

	category := c.DefaultQuery("category", "")

	products, total, err := h.getProductsUseCase.Run(c.Request.Context(), &ports.ProductFilter{
		Category: category,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": mappers.ToProductsDTO(products),
		"total":    total,
	})
}
