package models

import "time"

type Product struct {
	ID          string           `bson:"_id" json:"id"`
	TenantID    string           `bson:"tenant_id" json:"tenant_id"`
	PageID      string           `bson:"page_id" json:"page_id"`
	Name        string           `bson:"name" json:"name"`
	Description string           `bson:"description,omitempty" json:"description,omitempty"`
	Type        string           `bson:"type,omitempty" json:"type,omitempty"`
	PriceCents  int64            `bson:"price_cents" json:"price_cents"`
	Images      []string         `bson:"images,omitempty" json:"images,omitempty"`
	Visible     bool             `bson:"visible" json:"visible"`
	Variants    []ProductVariant `bson:"variants,omitempty" json:"variants,omitempty"`
	CreatedAt   time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time        `bson:"updated_at" json:"updated_at"`
}

type ProductVariant struct {
	Label      string `bson:"label,omitempty" json:"label,omitempty"`
	Model      string `bson:"model,omitempty" json:"model,omitempty"`
	SKU        string `bson:"sku,omitempty" json:"sku,omitempty"`
	PriceCents int64  `bson:"price_cents,omitempty" json:"price_cents,omitempty"`
	Quantity   int64  `bson:"quantity,omitempty" json:"quantity,omitempty"`
}
