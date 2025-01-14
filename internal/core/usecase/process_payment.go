package usecase

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type ProcessPaymentUseCase interface {
	Run(ctx context.Context, ExternalReference string, paymentMethod string, status entities.PaymentStatus) error
}

type processPaymentUseCase struct {
	paymentRepository ports.PaymentRepository
}

func NewProcessPaymentUseCase(paymentRepository ports.PaymentRepository) ProcessPaymentUseCase {
	return &processPaymentUseCase{paymentRepository: paymentRepository}
}

func (p *processPaymentUseCase) Run(ctx context.Context, ExternalReference string, paymentMethod string, status entities.PaymentStatus) error {
	return p.paymentRepository.UpdateOrderPaymentStatus(ctx, ExternalReference, paymentMethod, status)
}
