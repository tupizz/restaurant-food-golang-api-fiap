package usecase

import (
	"context"
	"errors"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type UpdateOrderStatusToReadyUseCase interface {
	Run(ctx context.Context, id int) error
}

type updateOrderStatusToReadyUseCase struct {
	orderRepository ports.OrderRepository
}

func NewUpdateOrderStatusToReadyUseCase(orderRepository ports.OrderRepository) UpdateOrderStatusToReadyUseCase {
	return &updateOrderStatusToReadyUseCase{orderRepository: orderRepository}
}

func (c *updateOrderStatusToReadyUseCase) Run(ctx context.Context, id int) error {
	order, err := c.orderRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != entities.OrderStatusPreparing {
		return errors.New("order status is not preparing")
	}

	return c.orderRepository.UpdateStatus(ctx, id, "ready")
}
