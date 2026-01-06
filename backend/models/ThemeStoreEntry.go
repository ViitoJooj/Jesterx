package models

import "time"

type ThemeStoreEntry struct {
	ID        string    `bson:"_id"`
	TenantID  string    `bson:"tenant_id"`
	PageID    string    `bson:"page_id"`
	Name      string    `bson:"name"`
	Domain    string    `bson:"domain,omitempty"`
	ForSale   bool      `bson:"for_sale"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
