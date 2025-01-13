package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	domainError "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/error"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/validator"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/shared"

	"github.com/gin-gonic/gin"
)

type ProductAdminHandler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type productAdminHandler struct {
	productUseCase usecase.ProductUseCase
}

func NewProductAdminHandler(productUseCase usecase.ProductUseCase) ProductAdminHandler {
	return &productAdminHandler{productUseCase: productUseCase}
}

// Create godoc
// @Summary      Create Product
// @Description  Create Product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        input  body     dto.ProductInputCreate  true  "Product data"
// @Success      201  {object}  dto.ProductOutput
// @Failure      400  {object}  handler.ErrorResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /admin/products [post]
func (h *productAdminHandler) Create(c *gin.Context) {
	var input dto.ProductInputCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.ValidateProductCreate(input); err != nil {
		errors := validator.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	slog.Info("Creating product", "input", shared.ToJSON(input))

	product, err := h.productUseCase.CreateProduct(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, &domainError.NotFoundError{}) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Product created", "product", shared.ToJSON(product))
	c.JSON(http.StatusCreated, product)
}

// Update godoc
// @Summary      Update Product
// @Description  Update Product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id     path     int  true  "Product ID"
// @Param        input  body     dto.ProductInputUpdate  true  "Product data"
// @Success      200  {object}  dto.ProductOutput
// @Failure      400  {object}  handler.ErrorResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /admin/products/{id} [put]
func (h *productAdminHandler) Update(c *gin.Context) {
	var input dto.ProductInputUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.ValidateProductUpdate(input); err != nil {
		errors := validator.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	slog.Info("Updating product", "input", shared.ToJSON(input))

	product, err := h.productUseCase.UpdateProduct(c.Request.Context(), input.ID, input)
	if err != nil {
		if errors.Is(err, &domainError.NotFoundError{}) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Product updated", "product", shared.ToJSON(product))
	c.JSON(http.StatusOK, product)
}

// Delete godoc
// @Summary      Delete Product
// @Description  Delete Product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id     path     int  true  "Product ID"
// @Success      204  {object}  dto.ProductOutput
// @Failure      400  {object}  handler.ErrorResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /admin/products/{id} [delete]
func (h *productAdminHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	slog.Info("Deleting product", "id", id)

	err = h.productUseCase.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Product deleted", "id", id)

	c.Status(http.StatusNoContent)
}
