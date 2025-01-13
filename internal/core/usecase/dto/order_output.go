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

type OrderDTO struct {
	ID       int            `json:"id"`
	ClientID int            `json:"client_id"`
	Client   ClientDTO      `json:"client"`
	Status   string         `json:"status"`
	Items    []OrderItemDTO `json:"items"`
	Payment  PaymentDTO     `json:"payment"`
}

type ClientDTO struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type OrderItemDTO struct {
	ID        int        `json:"id"`
	OrderID   int        `json:"order_id"`
	ProductID int        `json:"product_id"`
	Product   ProductDTO `json:"product"`
	Quantity  int        `json:"quantity"`
	Price     float64    `json:"price"`
}

type ProductDTO struct {
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Price          float64     `json:"price"`
	CategoryHandle string      `json:"category_handle"`
	Images         interface{} `json:"images,omitempty"`
	CreatedAt      time.Time   `json:"-"`
	UpdatedAt      time.Time   `json:"-"`
}

type PaymentDTO struct {
	ID      int     `json:"id"`
	OrderID int     `json:"order_id"`
	Status  string  `json:"status"`
	Method  string  `json:"method"`
	Amount  float64 `json:"amount"`
}

type PaginatedOrdersDTO struct {
	Orders []OrderDTO `json:"orders"`
	Total  int        `json:"total"`
}
