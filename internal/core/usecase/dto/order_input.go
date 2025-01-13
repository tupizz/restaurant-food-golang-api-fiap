package dto

type CreateOrderRequest struct {
	ClientID int                      `json:"client_id"`
	Items    []CreateOrderItemRequest `json:"items"`
	Payment  CreatePaymentRequest     `json:"payment"`
}

type CreateOrderItemRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type CreatePaymentRequest struct {
	Method string `json:"method" binding:"required"`
}
