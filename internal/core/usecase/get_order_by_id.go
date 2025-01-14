package usecase

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type GetOrderByIDUseCase interface {
	Run(ctx context.Context, id int) (*entities.Order, error)
}

type getOrderByIDUseCase struct {
	orderRepository ports.OrderRepository
}

func NewGetOrderByIDUseCase(orderRepository ports.OrderRepository) GetOrderByIDUseCase {
	return &getOrderByIDUseCase{orderRepository: orderRepository}
}

func (s *getOrderByIDUseCase) Run(ctx context.Context, id int) (*entities.Order, error) {
	order, err := s.orderRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
