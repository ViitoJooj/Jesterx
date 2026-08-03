package dtos

type CreateWebsiteRequest struct {
	OwnerUUID   string `json:"owner_uuid"`
	OwnerType   string `json:"owner_type"`
	Label       string `json:"label"`
	URL         string `json:"url"`
	WriteIn     string `json:"write_in"`
	Description string `json:"description"`
}

type WebsiteResponse struct {
	UUID        string `json:"uuid"`
	OwnerUUID   string `json:"owner_uuid"`
	OwnerType   string `json:"owner_type"`
	Label       string `json:"label"`
	URL         string `json:"url"`
	WriteIn     string `json:"write_in"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}
