package gateways

import (
	"github.com/google/uuid"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
)

type creditCardMock struct{}

func NewCreditCardMockGateway() PaymentGateway {
	return &creditCardMock{}
}

func (g *creditCardMock) Authorize(payment *entities.Payment) error {
	payment.ExternalReference = uuid.New().String()

	return nil
}
