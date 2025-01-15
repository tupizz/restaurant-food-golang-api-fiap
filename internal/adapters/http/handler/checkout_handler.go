package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainError "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/error"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/mappers"
)

type CheckoutHandler interface {
	Create(c *gin.Context)
}

type checkoutHandler struct {
	createOrderUseCase usecase.CreateOrderUseCase
}

func NewCheckoutHandler(createOrderUseCase usecase.CreateOrderUseCase) CheckoutHandler {
	return &checkoutHandler{createOrderUseCase: createOrderUseCase}
}

// Create godoc
// @Summary      Cria um novo pedido
// @Description  Cria um novo pedido com os dados fornecidos
// @Tags         checkout
// @Accept       json
// @Produce      json
// @Param        order  body      dto.CreateOrderRequest  true  "Dados do Pedido"
// @Success      201     {object}  dto.OrderResponse
// @Failure      400     {object}  handler.ErrorResponse
// @Failure      500     {object}  handler.ErrorResponse
// @Router       /checkout [post]
func (h *checkoutHandler) Create(c *gin.Context) {
	var createOrderReq dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&createOrderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderEntity := mappers.MapCreateOrderRequestToEntity(createOrderReq)

	createdOrder, err := h.createOrderUseCase.Run(c.Request.Context(), orderEntity)
	if err != nil {
		if errors.Is(err, &domainError.EntityNotProcessableError{}) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, mappers.MapOrderEntityToResponse(*createdOrder))
}
