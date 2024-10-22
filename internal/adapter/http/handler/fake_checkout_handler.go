package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto/payment_dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/service"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/validation"
)

type FakeCheckoutHandler interface {
	ProcessPayment(c *gin.Context)
}

type fakeCheckoutHandler struct {
	paymentService service.PaymentService
}

func NewFakeCheckoutHandler(paymentService service.PaymentService) FakeCheckoutHandler {
	return &fakeCheckoutHandler{paymentService: paymentService}
}

func (h fakeCheckoutHandler) ProcessPayment(c *gin.Context) {
	var paymentInput payment_dto.PaymentInputDto

	if err := c.ShouldBindJSON(&paymentInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate paymentInput
	if err := validator.New().Struct(paymentInput); err != nil {
		errors := validation.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	err := h.paymentService.ProcessPayment(c.Request.Context(), paymentInput.OrderId, paymentInput.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
}
