package entities

import (
	"fmt"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
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

// CalculateTotalAmount calculates the total amount for the order, including product prices and applicable taxes.
// It takes two parameters:
// - existingMappedProducts: a map of product IDs to Product entities
// - existingPaymentTaxes: a slice of PaymentTaxSettings entities
//
// The function performs the following steps:
// 1. Calculates the base total amount from the order items and their quantities
// 2. Applies any applicable taxes based on the payment method and tax settings
// 3. Sets the final amount to the Payment.Amount field of the Order
//
// Returns an error if any product in the order is not found in the existingMappedProducts map.
func (o *Order) CalculateTotalAmount(existingMappedProducts map[int]Product) error {
	totalAmount := 0.0

	// Calculate base total amount from order items
	for idx, item := range o.Items {
		if product, ok := existingMappedProducts[item.ProductID]; ok {
			item.Price = product.Price
			totalAmount += item.Price * float64(item.Quantity)
		} else {
			return fmt.Errorf("product not found for id %d", item.ProductID)
		}
		o.Items[idx] = item
	}

	o.Payment.Amount = totalAmount

	return nil
}
