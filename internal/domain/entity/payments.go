package entity

import "time"

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusApproved PaymentStatus = "approved"
	PaymentStatusFailed   PaymentStatus = "failed"
)

type PaymentMethod string

const (
	PaymentMethodPix        PaymentMethod = "pix"
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodBillet     PaymentMethod = "billet"
)

type Payment struct {
	ID        int
	OrderID   int
	Status    PaymentStatus
	Method    PaymentMethod
	Amount    float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
