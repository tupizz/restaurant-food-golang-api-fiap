package dto

import "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"

type PaymentInputDTO struct {
	ExternalReference string               `json:"external_reference" validate:"required"`
	Status            domain.PaymentStatus `json:"status" validate:"required,oneof=approved failed"`
	PaymentMethod     string               `json:"payment_method" validate:"required,oneof=qr_code"`
}
