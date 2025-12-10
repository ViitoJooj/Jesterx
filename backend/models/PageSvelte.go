package models

import "time"

type PageSvelte struct {
	ID        string    `bson:"_id"`
	TenantID  string    `bson:"tenant_id"`
	PageID    string    `bson:"page_id"`
	Svelte    string    `bson:"svelte"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
