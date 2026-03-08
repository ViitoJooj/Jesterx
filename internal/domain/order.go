package domain

import "time"

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderPaid      OrderStatus = "paid"
	OrderShipped   OrderStatus = "shipped"
	OrderDelivered OrderStatus = "delivered"
	OrderCancelled OrderStatus = "cancelled"
	OrderRefunded  OrderStatus = "refunded"
)

type OrderItem struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"order_id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	UnitPrice   float64 `json:"unit_price"`
	Qty         int     `json:"qty"`
	Total       float64 `json:"total"`
}

type Order struct {
	ID          string      `json:"id"`
	WebsiteID   string      `json:"website_id"`
	BuyerName   string      `json:"buyer_name"`
	BuyerEmail  string      `json:"buyer_email"`
	BuyerPhone  string      `json:"buyer_phone,omitempty"`
	Status      OrderStatus `json:"status"`
	Subtotal    float64     `json:"subtotal"`
	PlatformFee float64     `json:"platform_fee"`
	Total       float64     `json:"total"`
	Notes       string      `json:"notes,omitempty"`
	Items       []OrderItem `json:"items"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
