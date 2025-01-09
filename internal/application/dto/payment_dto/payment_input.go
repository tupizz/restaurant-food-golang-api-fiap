package payment_dto

import "github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"

type PaymentInputDto struct {
	ExternalReference string               `json:"external_reference" validate:"required"`
	Status            entity.PaymentStatus `json:"status" validate:"required,oneof=approved failed"`
	PaymentMethod     string               `json:"payment_method" validate:"required,oneof=qr_code"`
}
