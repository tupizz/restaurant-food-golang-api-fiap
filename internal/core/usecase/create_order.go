package usecase

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	domainError "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/error"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type CreateOrderUseCase interface {
	Run(ctx context.Context, order entities.Order) (*entities.Order, error)
}

type createOrderUseCase struct {
	orderRepository              ports.OrderRepository
	productRepository            ports.ProductRepository
	paymentTaxSettingsRepository ports.PaymentTaxSettingsRepository
}

func NewCreateOrderUseCase(orderRepository ports.OrderRepository, productRepository ports.ProductRepository, paymentTaxSettingsRepository ports.PaymentTaxSettingsRepository) CreateOrderUseCase {
	return &createOrderUseCase{orderRepository: orderRepository, productRepository: productRepository, paymentTaxSettingsRepository: paymentTaxSettingsRepository}
}

func (c *createOrderUseCase) Run(ctx context.Context, order entities.Order) (*entities.Order, error) {
	var productIds []int
	for _, item := range order.Items {
		productIds = append(productIds, item.ProductID)
	}

	existingProductsFromDB, _, err := c.productRepository.GetByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	mappedProducts := make(map[int]entities.Product)
	for _, product := range existingProductsFromDB {
		mappedProducts[product.ID] = product
	}

	systemPaymentTaxSettings, err := c.paymentTaxSettingsRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	err = order.CalculateTotalAmount(mappedProducts, systemPaymentTaxSettings)
	if err != nil {
		return nil, domainError.NewEntityNotProcessableError("order", err.Error())
	}

	err = order.Payment.Authorize()
	if err != nil {
		return nil, domainError.NewEntityNotProcessableError("payment", err.Error())
	}

	createdOrder, err := c.orderRepository.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return &createdOrder, nil
}
