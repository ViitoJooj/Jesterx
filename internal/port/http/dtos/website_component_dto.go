package dtos

import "encoding/json"

type CreateWebsiteComponentRequest struct {
	WebsiteUUID string          `json:"website_uuid"`
	LogoURL     string          `json:"logo_url"`
	Tittle      string          `json:"tittle"`
	Description string          `json:"description"`
	Path        string          `json:"path"`
	Content     json.RawMessage `json:"content"`
	Visits      int             `json:"visits"`
}

type WebsiteComponentResponse struct {
	UUID        string          `json:"uuid"`
	WebsiteUUID string          `json:"website_uuid"`
	LogoURL     string          `json:"logo_url"`
	Tittle      string          `json:"tittle"`
	Description string          `json:"description"`
	Path        string          `json:"path"`
	Content     json.RawMessage `json:"content"`
	Visits      int             `json:"visits"`
	CreatedAt   string          `json:"created_at"`
}
