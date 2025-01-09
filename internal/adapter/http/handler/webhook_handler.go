package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto/payment_dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/service"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/validation"
)

type WebhookHandler interface {
	ProcessPayment(c *gin.Context)
}

type webhookHandler struct {
	paymentService service.PaymentService
}

func NewWebhookHandler(paymentService service.PaymentService) WebhookHandler {
	return &webhookHandler{paymentService: paymentService}
}

func (h webhookHandler) ProcessPayment(c *gin.Context) {
	var paymentInput payment_dto.PaymentInputDto

	if err := c.ShouldBindJSON(&paymentInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(paymentInput); err != nil {
		errors := validation.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	err := h.paymentService.ProcessPayment(c.Request.Context(), paymentInput.ExternalReference, paymentInput.PaymentMethod, paymentInput.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
}
