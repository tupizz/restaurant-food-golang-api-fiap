package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapter/repository"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/service"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/mappers"
)

type OrderHandler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
}

type orderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) OrderHandler {
	return &orderHandler{orderService: orderService}
}

// CreateOrder handles POST /pedidos
func (h *orderHandler) Create(c *gin.Context) {
	var createOrderReq dto.CreateOrderRequest

	// Parse and validate the JSON body into the CreateOrderRequest DTO
	if err := c.ShouldBindJSON(&createOrderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map DTO to domain entity
	orderEntity := mappers.MapCreateOrderRequestToEntity(createOrderReq)

	// Call the service to create the order
	createdOrder, err := h.orderService.CreateOrder(c.Request.Context(), orderEntity)
	if err != nil {
		if errors.Is(err, &domain.EntityNotProcessableError{}) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Map domain entity to response DTO
	orderResponse := mappers.MapOrderEntityToResponse(createdOrder)

	// Return the created order as JSON
	c.JSON(http.StatusCreated, orderResponse)
}

// GetOrderById handles GET /pedidos/{id}
func (h *orderHandler) GetById(c *gin.Context) {
	// Get the order ID from the URL path
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Call the service to retrieve the order by ID
	orderEntity, err := h.orderService.GetOrderById(c.Request.Context(), orderID)
	if err != nil {
		if err == repository.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Map domain entity to response DTO
	orderResponse := mappers.MapOrderEntityToResponse(orderEntity)

	// Return the order as JSON
	c.JSON(http.StatusOK, orderResponse)
}
