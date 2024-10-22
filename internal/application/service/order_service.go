package service

import (
	"context"
	"log/slog"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto/order_list_dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error)
	GetOrderById(ctx context.Context, id int) (entity.Order, error)
	GetAllOrders(ctx context.Context, filter *domain.OrderFilter) ([]order_list_dto.OrderDTO, error)
}

type orderService struct {
	orderRepo              domain.OrderRepository
	productRepo            domain.ProductRepository
	paymentTaxSettingsRepo domain.PaymentTaxSettingsRepository
}

func NewOrderService(
	orderRepo domain.OrderRepository,
	productRepo domain.ProductRepository,
	paymentTaxSettingsRepo domain.PaymentTaxSettingsRepository,
) OrderService {
	return &orderService{
		orderRepo:              orderRepo,
		productRepo:            productRepo,
		paymentTaxSettingsRepo: paymentTaxSettingsRepo,
	}
}

func (s *orderService) GetAllOrders(ctx context.Context, filter *domain.OrderFilter) ([]order_list_dto.OrderDTO, error) {
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

func (s *orderService) CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error) {
	// Validate if the product exists
	var productIds []int
	for _, item := range order.Items {
		productIds = append(productIds, item.ProductID)
	}

	existingProductsFromDB, _, err := s.productRepo.GetByIds(ctx, productIds)
	if err != nil {
		return entity.Order{}, err
	}

	mappedProducts := make(map[int]entity.Product)
	for _, product := range existingProductsFromDB {
		mappedProducts[product.ID] = product
	}

	systemPaymentTaxSettings, err := s.paymentTaxSettingsRepo.GetAll(ctx)
	if err != nil {
		return entity.Order{}, err
	}

	err = order.CalculateTotalAmount(mappedProducts, systemPaymentTaxSettings)
	if err != nil {
		return entity.Order{}, domain.NewEntityNotProcessableError("order", err.Error())
	}

	// Business logic to create an order
	createdOrder, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return entity.Order{}, err
	}
	return createdOrder, nil
}

func (s *orderService) GetOrderById(ctx context.Context, id int) (entity.Order, error) {
	// Business logic to retrieve an order by ID
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return entity.Order{}, err
	}
	return order, nil
}
