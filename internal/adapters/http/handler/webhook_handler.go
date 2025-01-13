package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	ineternalValidator "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/validator"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
)

type WebhookHandler interface {
	ProcessPayment(c *gin.Context)
}

type webhookHandler struct {
	paymentUseCase usecase.PaymentUseCase
}

func NewWebhookHandler(paymentUseCase usecase.PaymentUseCase) WebhookHandler {
	return &webhookHandler{paymentUseCase: paymentUseCase}
}

func (h *webhookHandler) ProcessPayment(c *gin.Context) {
	var paymentInput dto.PaymentInputDTO

	if err := c.ShouldBindJSON(&paymentInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(paymentInput); err != nil {
		errors := ineternalValidator.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	err := h.paymentUseCase.ProcessPayment(c.Request.Context(), paymentInput.ExternalReference, paymentInput.PaymentMethod, paymentInput.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
}
