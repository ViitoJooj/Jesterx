package dtos

type CreateAddressRequest struct {
	OwnerUUID      string `json:"owner_uuid"`
	OwnerType      string `json:"owner_type"`
	Label          string `json:"label"`
	AddressLine1   string `json:"address_line1"`
	AddressLine2   string `json:"address_line2"`
	Neighborhood   string `json:"neighborhood"`
	City           string `json:"city"`
	State          string `json:"state"`
	StateCode      string `json:"state_code"`
	PostalCode     string `json:"postal_code"`
	ReferencePoint string `json:"reference_point"`
	DeliveryNotes  string `json:"delivery_notes"`
	IsDefault      bool   `json:"is_default"`
}

type AddressResponse struct {
	UUID           string `json:"uuid"`
	WebSiteUUID    string `json:"website_uuid"`
	OwnerUUID      string `json:"owner_uuid"`
	OwnerType      string `json:"owner_type"`
	Label          string `json:"label"`
	AddressLine1   string `json:"address_line1"`
	AddressLine2   string `json:"address_line2"`
	Neighborhood   string `json:"neighborhood"`
	City           string `json:"city"`
	State          string `json:"state"`
	StateCode      string `json:"state_code"`
	PostalCode     string `json:"postal_code"`
	ReferencePoint string `json:"reference_point"`
	DeliveryNotes  string `json:"delivery_notes"`
	IsDefault      bool   `json:"is_default"`
	UpdatedAt      string `json:"updated_at"`
	CreatedAt      string `json:"created_at"`
}
