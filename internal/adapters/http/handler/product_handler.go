package handler

import (
	"net/http"
	"strconv"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"

	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	GetProducts(c *gin.Context)
}

type productHandler struct {
	productUseCase usecase.ProductUseCase
}

func NewProductHandler(productUseCase usecase.ProductUseCase) ProductHandler {
	return &productHandler{productUseCase: productUseCase}
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

	products, total, err := h.productUseCase.GetProducts(c.Request.Context(), &ports.ProductFilter{
		Category: category,
		Page:     page,
		PageSize: pageSize,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
	})
}