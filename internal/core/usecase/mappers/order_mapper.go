package mappers

import (
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
)

func MapCreateOrderRequestToEntity(dto dto.CreateOrderRequest) entities.Order {
	items := make([]entities.OrderItem, len(dto.Items))
	for i, itemDTO := range dto.Items {
		items[i] = entities.OrderItem{
			ProductID: itemDTO.ProductID,
			Quantity:  itemDTO.Quantity,
		}
	}

	payment := entities.Payment{
		Method: entities.PaymentMethod(dto.Payment.Method),
		Status: entities.PaymentStatusPending,
	}

	return entities.Order{
		ClientID: dto.ClientID,
		Status:   entities.OrderStatusPending,
		Items:    items,
		Payment:  payment,
	}
}

func MapOrderEntityToResponse(order entities.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dto.OrderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	payment := dto.PaymentResponse{
		ID:        order.Payment.ID,
		OrderID:   order.Payment.OrderID,
		Status:    string(order.Payment.Status),
		Method:    string(order.Payment.Method),
		QRData:    order.Payment.QRData,
		Amount:    order.Payment.Amount,
		CreatedAt: order.Payment.CreatedAt,
		UpdatedAt: order.Payment.UpdatedAt,
	}

	return dto.OrderResponse{
		ID:        order.ID,
		ClientID:  order.ClientID,
		Status:    string(order.Status),
		Items:     items,
		Payment:   payment,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}
