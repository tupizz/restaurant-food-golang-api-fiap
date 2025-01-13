package usecase

import (
	"context"
	"log/slog"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto/order_list_dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
	domainError "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/error"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type OrderUseCase interface {
	CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error)
	GetOrderById(ctx context.Context, id int) (domain.Order, error)
	GetAllOrders(ctx context.Context, filter *ports.OrderFilter) ([]order_list_dto.OrderDTO, error)
}

type orderUseCase struct {
	orderRepo              ports.OrderRepository
	productRepo            ports.ProductRepository
	paymentTaxSettingsRepo ports.PaymentTaxSettingsRepository
}

func NewOrderUseCase(
	orderRepo ports.OrderRepository,
	productRepo ports.ProductRepository,
	paymentTaxSettingsRepo ports.PaymentTaxSettingsRepository,
) OrderUseCase {
	return &orderUseCase{
		orderRepo:              orderRepo,
		productRepo:            productRepo,
		paymentTaxSettingsRepo: paymentTaxSettingsRepo,
	}
}

func (s *orderUseCase) GetAllOrders(ctx context.Context, filter *ports.OrderFilter) ([]order_list_dto.OrderDTO, error) {
	if filter.PageSize == 0 {
		filter.PageSize = 10
	}

	if filter.Page == 0 {
		filter.Page = 1
	}

	orders, err := s.orderRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	var mapOrderIdToItems = make(map[int][]order_list_dto.OrderDTO)
	for _, order := range orders {
		productPrice, _ := order.ProductPrice.Float64Value()

		if !productPrice.Valid {
			slog.Error("product price is not valid")
		}

		paymentAmount, _ := order.PaymentAmount.Float64Value()
		if !paymentAmount.Valid {
			slog.Error("payment amount is not valid")
		}

		mapOrderIdToItems[int(order.OrderID)] = append(mapOrderIdToItems[int(order.OrderID)], order_list_dto.OrderDTO{
			ID:       int(order.OrderID),
			ClientID: int(order.ClientID),
			Client: order_list_dto.ClientDTO{
				ID:   int(order.ClientID),
				Name: order.ClientName,
				CPF:  order.ClientCpf,
			},
			Status: string(order.OrderStatus.String),
			Items: []order_list_dto.OrderItemDTO{
				{
					ID:        int(order.ProductID),
					OrderID:   int(order.OrderID),
					ProductID: int(order.ProductID),
					Quantity:  int(order.ProductQuantity),
					Price:     productPrice.Float64,
					Product: order_list_dto.ProductDTO{
						ID:             int(order.ProductID),
						Name:           order.ProductName,
						CategoryHandle: order.CategoryHandle.String,
						Description:    order.ProductDescription,
						Price:          productPrice.Float64,
					},
				},
			},
			Payment: order_list_dto.PaymentDTO{
				ID:      int(order.PaymentID),
				OrderID: int(order.OrderID),
				Status:  string(order.PaymentStatus.String),
				Amount:  paymentAmount.Float64,
				Method:  string(order.PaymentMethod),
			},
		})
	}

	var ordersEntity []order_list_dto.OrderDTO
	for _, orders := range mapOrderIdToItems {
		order := orders[0]

		for _, item := range orders {
			order.Items = append(order.Items, item.Items...)
		}

		ordersEntity = append(ordersEntity, order)
	}

	return ordersEntity, nil
}

func (s *orderUseCase) CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	// Validate if the product exists
	var productIds []int
	for _, item := range order.Items {
		productIds = append(productIds, item.ProductID)
	}

	existingProductsFromDB, _, err := s.productRepo.GetByIds(ctx, productIds)
	if err != nil {
		return domain.Order{}, err
	}

	mappedProducts := make(map[int]domain.Product)
	for _, product := range existingProductsFromDB {
		mappedProducts[product.ID] = product
	}

	systemPaymentTaxSettings, err := s.paymentTaxSettingsRepo.GetAll(ctx)
	if err != nil {
		return domain.Order{}, err
	}

	err = order.CalculateTotalAmount(mappedProducts, systemPaymentTaxSettings)
	if err != nil {
		return domain.Order{}, domainError.NewEntityNotProcessableError("order", err.Error())
	}

	err = order.Payment.Authorize()
	if err != nil {
		return domain.Order{}, domainError.NewEntityNotProcessableError("payment", err.Error())
	}

	createdOrder, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return domain.Order{}, err
	}

	return createdOrder, nil
}

func (s *orderUseCase) GetOrderById(ctx context.Context, id int) (domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
