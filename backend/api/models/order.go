package models

import "time"

type Order struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	OrderNumber string      `json:"order_number"`
	Status      string      `json:"status"`
	Subtotal    float64     `json:"subtotal"`
	ShippingFee float64     `json:"shipping_fee"`
	Discount    float64     `json:"discount"`
	Total       float64     `json:"total"`
	AddressID   string      `json:"address_id"`
	Notes       string      `json:"notes,omitempty"`
	Items       []OrderItem `json:"items,omitempty"`
	Address     *Address    `json:"address,omitempty"`
	Payment     *Payment    `json:"payment,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID         string   `json:"id"`
	OrderID    string   `json:"order_id"`
	VariantID  string   `json:"variant_id"`
	Quantity   int      `json:"quantity"`
	UnitPrice  float64  `json:"unit_price"`
	TotalPrice float64  `json:"total_price"`
	Variant    *Variant `json:"variant,omitempty"`
}

// CreateOrderRequest is used when a user places a new order
type CreateOrderRequest struct {
	Items []struct {
		VariantID string `json:"variant_id" binding:"required"`
		Quantity  int    `json:"quantity"   binding:"required,min=1"`
	} `json:"items"      binding:"required,min=1"`
	AddressID string `json:"address_id" binding:"required"`
	Notes     string `json:"notes"`
}

// UpdateOrderStatusRequest is used by admin to update an order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed shipped delivered cancelled"`
}

// OrderFilterQuery holds query params for filtering orders
type OrderFilterQuery struct {
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
	Status string `form:"status"`
}
