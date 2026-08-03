package dtos

type CreateTermsRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TermsResponse struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}
