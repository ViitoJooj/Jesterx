package models

import "time"

type PageSvelte struct {
	ID         string    `bson:"_id"`
	TenantID   string    `bson:"tenant_id"`
	PageID     string    `bson:"page_id"`
	Svelte     string    `bson:"svelte"`
	Header     string    `bson:"header,omitempty"`
	Footer     string    `bson:"footer,omitempty"`
	ShowHeader bool      `bson:"show_header,omitempty"`
	ShowFooter bool      `bson:"show_footer,omitempty"`
	Components []string  `bson:"components,omitempty"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
}
