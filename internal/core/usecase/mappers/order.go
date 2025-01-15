package mappers

import (
	"log/slog"

	fiapRestaurantDb "github.com/tupizz/restaurant-food-golang-api-fiap/database/sqlc"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
)

func ToCompleteOrdersDTO(rawOrders []fiapRestaurantDb.GetAllOrdersRow) []dto.OrderDTO {
	var mapOrderIdToItems = make(map[int][]dto.OrderDTO)

	for _, order := range rawOrders {
		productPrice, _ := order.ProductPrice.Float64Value()

		if !productPrice.Valid {
			slog.Error("product price is not valid")
		}

		paymentAmount, _ := order.PaymentAmount.Float64Value()
		if !paymentAmount.Valid {
			slog.Error("payment amount is not valid")
		}

		mapOrderIdToItems[int(order.OrderID)] = append(mapOrderIdToItems[int(order.OrderID)], dto.OrderDTO{
			ID:       int(order.OrderID),
			ClientID: int(order.ClientID),
			Client: dto.ClientDTO{
				ID:   int(order.ClientID),
				Name: order.ClientName,
				CPF:  order.ClientCpf,
			},
			Status: string(order.OrderStatus.String),
			Items: []dto.OrderItemDTO{
				{
					ID:        int(order.ProductID),
					OrderID:   int(order.OrderID),
					ProductID: int(order.ProductID),
					Quantity:  int(order.ProductQuantity),
					Price:     productPrice.Float64,
					Product: dto.ProductDTO{
						ID:             int(order.ProductID),
						Name:           order.ProductName,
						CategoryHandle: order.CategoryHandle,
						Description:    order.ProductDescription,
						Price:          productPrice.Float64,
					},
				},
			},
			Payment: dto.PaymentDTO{
				ID:      int(order.PaymentID),
				OrderID: int(order.OrderID),
				Status:  string(order.PaymentStatus.String),
				Amount:  paymentAmount.Float64,
				Method:  string(order.PaymentMethod),
			},
		})
	}

	var ordersEntity []dto.OrderDTO
	for _, orders := range mapOrderIdToItems {
		order := orders[0]

		for _, item := range orders {
			order.Items = append(order.Items, item.Items...)
		}

		ordersEntity = append(ordersEntity, order)
	}

	return ordersEntity
}

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
