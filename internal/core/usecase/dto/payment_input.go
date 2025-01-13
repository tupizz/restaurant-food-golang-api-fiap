package dto

import "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"

type PaymentInputDTO struct {
	ExternalReference string                 `json:"external_reference" validate:"required"`
	Status            entities.PaymentStatus `json:"status" validate:"required,oneof=approved failed"`
	PaymentMethod     string                 `json:"payment_method" validate:"required,oneof=qr_code"`
}
