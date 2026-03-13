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
	ID                        string      `json:"id"`
	WebsiteID                 string      `json:"website_id"`
	BuyerUserID               *string     `json:"buyer_user_id,omitempty"`
	BuyerName                 string      `json:"buyer_name"`
	BuyerEmail                string      `json:"buyer_email"`
	BuyerPhone                string      `json:"buyer_phone,omitempty"`
	BuyerDocument             string      `json:"buyer_document,omitempty"`
	ShippingName              string      `json:"shipping_name,omitempty"`
	ShippingPhone             string      `json:"shipping_phone,omitempty"`
	ShippingZipCode           string      `json:"shipping_zip_code,omitempty"`
	ShippingAddressStreet     string      `json:"shipping_address_street,omitempty"`
	ShippingAddressNumber     string      `json:"shipping_address_number,omitempty"`
	ShippingAddressComplement string      `json:"shipping_address_complement,omitempty"`
	ShippingAddressDistrict   string      `json:"shipping_address_district,omitempty"`
	ShippingAddressCity       string      `json:"shipping_address_city,omitempty"`
	ShippingAddressState      string      `json:"shipping_address_state,omitempty"`
	ShippingAddressCountry    string      `json:"shipping_address_country,omitempty"`
	ShippingMethod            string      `json:"shipping_method,omitempty"`
	ShippingCost              float64     `json:"shipping_cost"`
	DiscountTotal             float64     `json:"discount_total"`
	TaxTotal                  float64     `json:"tax_total"`
	Currency                  string      `json:"currency"`
	Status                    OrderStatus `json:"status"`
	Subtotal                  float64     `json:"subtotal"`
	PlatformFee               float64     `json:"platform_fee"`
	Total                     float64     `json:"total"`
	Notes                     string      `json:"notes,omitempty"`
	Items                     []OrderItem `json:"items"`
	CreatedAt                 time.Time   `json:"created_at"`
	UpdatedAt                 time.Time   `json:"updated_at"`
}
