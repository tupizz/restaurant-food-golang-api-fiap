package payment_dto

import "github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"

type PaymentInputDto struct {
	OrderId int                  `json:"order_id" validate:"required"`
	Status  entity.PaymentStatus `json:"status" validate:"required,oneof=pending approved rejected"`
}
