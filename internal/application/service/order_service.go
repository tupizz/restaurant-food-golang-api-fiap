package service

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error)
	GetOrderById(ctx context.Context, id int) (entity.Order, error)
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
