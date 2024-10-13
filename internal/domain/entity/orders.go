package entity

import (
	"fmt"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusReceived  OrderStatus = "received"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCanceled  OrderStatus = "canceled"
)

type Order struct {
	ID        int
	ClientID  int
	Status    OrderStatus
	Items     []OrderItem
	Payment   Payment
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	Quantity  int
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (o *Order) CalculateTotalAmount(existingMappedProducts map[int]Product, existingPaymentTaxes []PaymentTaxSettings) error {
	totalAmount := 0.0

	for idx, item := range o.Items {
		if product, ok := existingMappedProducts[item.ProductID]; ok {
			item.Price = product.Price
			totalAmount += item.Price * float64(item.Quantity)
		} else {
			return fmt.Errorf("product not found for id %d", item.ProductID)
		}
		o.Items[idx] = item
	}

	for _, tax := range existingPaymentTaxes {
		if o.Payment.Method != "credit_card" && tax.ApplicableTo == ApplicableToCreditCard {
			continue
		}

		if tax.AmountType == AmountTypePercentage {
			totalAmount = totalAmount * (1 + tax.AmountValue/100)
		} else if tax.AmountType == AmountTypeFixed {
			totalAmount += tax.AmountValue
		}
	}

	o.Payment.Amount = totalAmount
	return nil
}
