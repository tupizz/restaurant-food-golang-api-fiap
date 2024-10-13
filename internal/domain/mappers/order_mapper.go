package mappers

import (
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

// MapCreateOrderRequestToEntity maps the CreateOrderRequest DTO to the Order entity
func MapCreateOrderRequestToEntity(dto dto.CreateOrderRequest) entity.Order {
	items := make([]entity.OrderItem, len(dto.Items))
	for i, itemDTO := range dto.Items {
		items[i] = entity.OrderItem{
			ProductID: itemDTO.ProductID,
			Quantity:  itemDTO.Quantity,
		}
	}

	payment := entity.Payment{
		Method: entity.PaymentMethod(dto.Payment.Method),
		Status: entity.PaymentStatusPending,
	}

	return entity.Order{
		ClientID: dto.ClientID,
		Status:   entity.OrderStatusPending,
		Items:    items,
		Payment:  payment,
	}
}

// MapOrderEntityToResponse maps the Order entity to the OrderResponse DTO
func MapOrderEntityToResponse(order entity.Order) dto.OrderResponse {
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
