package usecase

import (
	"context"

	fiapRestaurantDb "github.com/tupizz/restaurant-food-golang-api-fiap/database/sqlc"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type GetAllOrdersUseCase interface {
	Run(ctx context.Context, filter *ports.OrderFilter) ([]fiapRestaurantDb.GetAllOrdersRow, error)
}

type getAllOrdersUseCase struct {
	orderRepository ports.OrderRepository
}

func NewOrderUseCase(orderRepository ports.OrderRepository) GetAllOrdersUseCase {
	return &getAllOrdersUseCase{orderRepository: orderRepository}
}

func (s *getAllOrdersUseCase) Run(ctx context.Context, filter *ports.OrderFilter) ([]fiapRestaurantDb.GetAllOrdersRow, error) {
	if filter.PageSize == 0 {
		filter.PageSize = 10
	}

	if filter.Page == 0 {
		filter.Page = 1
	}

	orders, err := s.orderRepository.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
