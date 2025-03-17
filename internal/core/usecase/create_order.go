package usecase

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	paymentGateways "github.com/tupizz/restaurant-food-golang-api-fiap/internal/adapters/gateways/payment"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	domainError "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/error"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type CreateOrderUseCase interface {
	Run(ctx context.Context, order entities.Order) (*entities.Order, error)
}

type createOrderUseCase struct {
	orderRepository   ports.OrderRepository
	productRepository ports.ProductRepository
	redisClient       *redis.Client
}

func NewCreateOrderUseCase(orderRepository ports.OrderRepository, productRepository ports.ProductRepository, redisClient *redis.Client) CreateOrderUseCase {
	return &createOrderUseCase{orderRepository: orderRepository, productRepository: productRepository, redisClient: redisClient}
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

	err = order.CalculateTotalAmount(mappedProducts)
	if err != nil {
		return nil, domainError.NewEntityNotProcessableError("order", err.Error())
	}

	var paymentGateway ports.PaymentGateway
	switch order.Payment.Method {
	case "qr_code":
		paymentGateway = paymentGateways.NewQRCodePaymentGateway(c.redisClient)
	case "credit_card":
		paymentGateway = paymentGateways.NewCreditCardMockGateway()
	default:
		return nil, errors.New("payment method not supported")
	}

	err = paymentGateway.Authorize(&order.Payment)
	if err != nil {
		return nil, domainError.NewEntityNotProcessableError("payment", err.Error())
	}

	createdOrder, err := c.orderRepository.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return &createdOrder, nil
}
