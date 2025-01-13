package dto

import "time"

type OrderResponse struct {
	ID        int                 `json:"id"`
	ClientID  int                 `json:"client_id"`
	Status    string              `json:"status"`
	Items     []OrderItemResponse `json:"items"`
	Payment   PaymentResponse     `json:"payment"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type OrderItemResponse struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaymentResponse struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	Status    string    `json:"status"`
	Method    string    `json:"method"`
	QRData    string    `json:"qr_data,omitempty"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
