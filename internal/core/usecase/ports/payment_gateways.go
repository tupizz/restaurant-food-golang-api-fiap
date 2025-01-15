package ports

import "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"

type PaymentGateway interface {
	Authorize(payment *entities.Payment) error
}
