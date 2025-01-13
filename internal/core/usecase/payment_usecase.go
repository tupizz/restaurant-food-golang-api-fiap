package usecase

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, ExternalReference string, paymentMethod string, status domain.PaymentStatus) error
}

type paymentUseCase struct {
	paymentRepository ports.PaymentRepository
}

func NewPaymentUseCase(paymentRepository ports.PaymentRepository) PaymentUseCase {
	return &paymentUseCase{paymentRepository: paymentRepository}
}

func (p *paymentUseCase) ProcessPayment(ctx context.Context, ExternalReference string, paymentMethod string, status domain.PaymentStatus) error {
	return p.paymentRepository.UpdateOrderPaymentStatus(ctx, ExternalReference, paymentMethod, status)
}
