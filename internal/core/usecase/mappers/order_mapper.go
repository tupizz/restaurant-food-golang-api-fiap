package mappers

import (
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
)

func MapCreateOrderRequestToDomain(dto dto.CreateOrderRequest) domain.Order {
	items := make([]domain.OrderItem, len(dto.Items))
	for i, itemDTO := range dto.Items {
		items[i] = domain.OrderItem{
			ProductID: itemDTO.ProductID,
			Quantity:  itemDTO.Quantity,
		}
	}

	payment := domain.Payment{
		Method: domain.PaymentMethod(dto.Payment.Method),
		Status: domain.PaymentStatusPending,
	}

	return domain.Order{
		ClientID: dto.ClientID,
		Status:   domain.OrderStatusPending,
		Items:    items,
		Payment:  payment,
	}
}

func MapOrderDomainToResponse(order domain.Order) dto.OrderResponse {
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
