package model

import "time"

// Order represents a purchase order.
type Order struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customer_id"`
	Status      string      `json:"status"`
	TotalAmount float64     `json:"total_amount"`
	Items       []OrderItem `json:"items,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// CreateOrderRequest is the payload for POST /orders.
type CreateOrderRequest struct {
	CustomerID string           `json:"customer_id"`
	Items      []OrderItemInput `json:"items"`
}

// OrderItemInput describes one line in a CreateOrderRequest.
type OrderItemInput struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
