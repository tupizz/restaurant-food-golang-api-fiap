package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

const lockTTL = 10 * time.Second

type ProcessPaymentUseCase interface {
	Run(ctx context.Context, ExternalReference string, paymentMethod string, status entities.PaymentStatus) error
}

type processPaymentUseCase struct {
	paymentRepository ports.PaymentRepository
	redisClient       *redis.Client
}

func NewProcessPaymentUseCase(paymentRepository ports.PaymentRepository, redisClient *redis.Client) ProcessPaymentUseCase {
	return &processPaymentUseCase{
		paymentRepository: paymentRepository,
		redisClient:       redisClient,
	}
}

func (p *processPaymentUseCase) Run(ctx context.Context, ExternalReference string, paymentMethod string, status entities.PaymentStatus) error {
	lockKey := fmt.Sprintf("lock:payment:%s", ExternalReference)

	locked, err := p.redisClient.SetNX(ctx, lockKey, 1, lockTTL).Result()
	if err != nil {
		return fmt.Errorf("error acquiring lock: %w", err)
	} else if !locked {
		return fmt.Errorf("payment processing is already in progress for %s", ExternalReference)
	}
	defer p.redisClient.Del(ctx, lockKey)

	return p.paymentRepository.UpdateOrderPaymentStatus(ctx, ExternalReference, paymentMethod, status)
}
