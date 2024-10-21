package order_list

import "time"

// OrderDTO represents the top-level order response
type OrderDTO struct {
	ID       int            `json:"id"`
	ClientID int            `json:"client_id"`
	Client   ClientDTO      `json:"client"`
	Status   string         `json:"status"`
	Items    []OrderItemDTO `json:"items"`
	Payment  PaymentDTO     `json:"payment"`
}

// ClientDTO represents client information within an order
type ClientDTO struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// OrderItemDTO represents an individual item in an order
type OrderItemDTO struct {
	ID        int        `json:"id"`
	OrderID   int        `json:"order_id"`
	ProductID int        `json:"product_id"`
	Product   ProductDTO `json:"product"`
	Quantity  int        `json:"quantity"`
	Price     float64    `json:"price"`
}

// ProductDTO represents product information for each order item
type ProductDTO struct {
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Price          float64     `json:"price"`
	CategoryHandle string      `json:"category_handle"`
	Images         interface{} `json:"images,omitempty"` // Change this according to your image structure
	CreatedAt      time.Time   `json:"-"`
	UpdatedAt      time.Time   `json:"-"`
}

// PaymentDTO represents payment information for an order
type PaymentDTO struct {
	ID      int     `json:"id"`
	OrderID int     `json:"order_id"`
	Status  string  `json:"status"`
	Method  string  `json:"method"`
	Amount  float64 `json:"amount"`
}

// PaginatedOrdersDTO represents the paginated response
type PaginatedOrdersDTO struct {
	Orders []OrderDTO `json:"orders"`
	Total  int        `json:"total"`
}
