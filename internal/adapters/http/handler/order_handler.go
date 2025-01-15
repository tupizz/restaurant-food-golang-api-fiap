package handler

import (
	"net/http"
	"strconv"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapters/db/repository"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/mappers"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"

	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	GetById(c *gin.Context)
	GetAll(c *gin.Context)
	UpdateOrderStatusToReady(c *gin.Context)
	UpdateOrderStatusToDelivered(c *gin.Context)
}

type orderHandler struct {
	getAllOrdersUseCase                 usecase.GetAllOrdersUseCase
	getOrderByIDUseCase                 usecase.GetOrderByIDUseCase
	updateOrderStatusToReadyUseCase     usecase.UpdateOrderStatusToReadyUseCase
	updateOrderStatusToDeliveredUseCase usecase.UpdateOrderStatusToDeliveredUseCase
}

func NewOrderHandler(getAllOrdersUseCase usecase.GetAllOrdersUseCase, getOrderByIDUseCase usecase.GetOrderByIDUseCase, updateOrderStatusToReadyUseCase usecase.UpdateOrderStatusToReadyUseCase, updateOrderStatusToDeliveredUseCase usecase.UpdateOrderStatusToDeliveredUseCase) OrderHandler {
	return &orderHandler{getAllOrdersUseCase: getAllOrdersUseCase, getOrderByIDUseCase: getOrderByIDUseCase, updateOrderStatusToReadyUseCase: updateOrderStatusToReadyUseCase, updateOrderStatusToDeliveredUseCase: updateOrderStatusToDeliveredUseCase}
}

// GetById godoc
// @Summary      Obtém um pedido por ID
// @Description  Obtém um pedido com o ID fornecido
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "ID do Pedido"
// @Success      200     {object}  dto.OrderResponse
// @Failure      400     {object}  handler.ErrorResponse
// @Failure      500     {object}  handler.ErrorResponse
// @Router       /orders/{id} [get]
func (h *orderHandler) GetById(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	orderEntity, err := h.getOrderByIDUseCase.Run(c.Request.Context(), orderID)
	if err != nil {
		if err == repository.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, mappers.MapOrderEntityToResponse(*orderEntity))
}

// GetAllOrders godoc
// @Summary     Retrieve all orders
// @Description Get a list of all orders with pagination
// @Tags        orders
// @Accept      json
// @Produce     json
// @Param       page     query     int    false  "Page number"        default(1)
// @Param       pageSize query     int    false  "Number of items per page" default(10)
// @Success     200      {object}  dto.PaginatedOrdersDTO
// @Failure     400      {object}  ErrorResponse
// @Failure     500      {object}  ErrorResponse
// @Router      /admin/orders [get]
func (h *orderHandler) GetAll(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
	}

	orders, err := h.getAllOrdersUseCase.Run(c.Request.Context(), &ports.OrderFilter{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": mappers.ToCompleteOrdersDTO(orders),
		"total":  len(orders),
	})
}

// UpdateOrderStatusToReady godoc
// @Summary     Mark order as ready
// @Description Mark ordeer as ready
// @Tags        orders
// @Accept      json
// @Produce     json
// @Param       id path     int    true  "Order ID"
// @Success      204 "No content"
// @Failure     400      {object}  ErrorResponse
// @Router      /admin/orders/{id}/ready [patch]
func (h *orderHandler) UpdateOrderStatusToReady(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.updateOrderStatusToReadyUseCase.Run(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateOrderStatusToDelivered godoc
// @Summary     Mark order as delivered
// @Description Mark ordeer as delivered
// @Tags        orders
// @Accept      json
// @Produce     json
// @Param       id path     int    true  "Order ID"
// @Success      204 "No content"
// @Failure     400      {object}  ErrorResponse
// @Router      /admin/orders/{id}/delivered [patch]
func (h *orderHandler) UpdateOrderStatusToDelivered(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.updateOrderStatusToDeliveredUseCase.Run(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
