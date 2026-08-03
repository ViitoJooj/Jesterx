package dtos

type CreateOrganizationRequest struct {
	OwnerUUID string `json:"owner_uuid"`
	ImageURL  string `json:"image_url"`
	Name      string `json:"name"`
	TradeName string `json:"trade_name"`
	CNPJ      string `json:"cnpj"`
}

type OrganizationResponse struct {
	UUID        string `json:"uuid"`
	WebSiteUUID string `json:"website_uuid"`
	OwnerUUID   string `json:"owner_uuid"`
	ImageURL    string `json:"image_url"`
	Name        string `json:"name"`
	TradeName   string `json:"trade_name"`
	CNPJ        string `json:"cnpj"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}
