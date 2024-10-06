package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/service"
)

type ProductHandler interface {
	GetProducts(c *gin.Context)
}

type productHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return &productHandler{productService: productService}
}

// GetProducts godoc
// @Summary      Get Products
// @Description  Get Products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page     query     int  false  "Page number"
// @Param        pageSize query     int  false  "Page size"
// @Success      200      {object}  dto.ProductOutput
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

	products, total, err := h.productService.GetProducts(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
	})
}
