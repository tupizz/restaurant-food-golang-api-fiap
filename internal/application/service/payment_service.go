package service

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, ExternalReference string, paymentMethod string, status entity.PaymentStatus) error
}

type paymentServiceImpl struct {
	paymentRepository domain.PaymentRepository
}

func NewPaymentService(paymentRepository domain.PaymentRepository) PaymentService {
	return &paymentServiceImpl{paymentRepository: paymentRepository}
}

func (s *paymentServiceImpl) ProcessPayment(ctx context.Context, ExternalReference string, paymentMethod string, status entity.PaymentStatus) error {
	return s.paymentRepository.UpdateOrderPaymentStatus(ctx, ExternalReference, paymentMethod, status)
}
