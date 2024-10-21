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
	GetAll(c *gin.Context)
}

type orderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) OrderHandler {
	return &orderHandler{orderService: orderService}
}

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

	orders, err := h.orderService.GetAllOrders(c.Request.Context(), &domain.OrderFilter{
		Page:     page,
		PageSize: pageSize,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  len(orders),
	})
}

// Create godoc
// @Summary      Cria um novo pedido
// @Description  Cria um novo pedido com os dados fornecidos
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order  body      dto.CreateOrderRequest  true  "Dados do Pedido"
// @Success      201     {object}  dto.OrderResponse
// @Failure      400     {object}  handler.ErrorResponse
// @Failure      500     {object}  handler.ErrorResponse
// @Router       /orders [post]
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
