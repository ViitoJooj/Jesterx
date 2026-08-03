package dtos

type CreatePhoneRequest struct {
	OwnerUUID string `json:"owner_uuid"`
	OwnerType string `json:"owner_type"`
	Label     string `json:"label"`
	Number    int    `json:"number"`
	IsDefault bool   `json:"is_default"`
}

type PhoneResponse struct {
	UUID        string `json:"uuid"`
	WebSiteUUID string `json:"website_uuid"`
	OwnerUUID   string `json:"owner_uuid"`
	OwnerType   string `json:"owner_type"`
	Label       string `json:"label"`
	Number      int    `json:"number"`
	IsDefault   bool   `json:"is_default"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}
