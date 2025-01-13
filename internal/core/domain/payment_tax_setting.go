package domain

import "time"

type AmountType string

const (
	AmountTypeFixed      AmountType = "fixed"
	AmountTypePercentage AmountType = "percentage"
)

type ApplicableTo string

const (
	ApplicableToCreditCard     ApplicableTo = "credit_card"
	ApplicableToTransportation ApplicableTo = "transportation"
	ApplicableToPlatformFee    ApplicableTo = "platform_fee"
)

type PaymentTaxSettings struct {
	ID           int
	Name         string
	Description  string
	AmountType   AmountType
	AmountValue  float64
	ApplicableTo ApplicableTo
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
