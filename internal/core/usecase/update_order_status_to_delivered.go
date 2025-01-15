package usecase

import (
	"context"
	"errors"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type UpdateOrderStatusToDeliveredUseCase interface {
	Run(ctx context.Context, id int) error
}

type updateOrderStatusToDeliveredUseCase struct {
	orderRepository ports.OrderRepository
}

func NewUpdateOrderStatusToDeliveredUseCase(orderRepository ports.OrderRepository) UpdateOrderStatusToDeliveredUseCase {
	return &updateOrderStatusToDeliveredUseCase{orderRepository: orderRepository}
}

func (c *updateOrderStatusToDeliveredUseCase) Run(ctx context.Context, id int) error {
	order, err := c.orderRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != entities.OrderStatusReady {
		return errors.New("order status is not ready")
	}

	return c.orderRepository.UpdateStatus(ctx, id, "delivered")
}
