package dtos

type CreateProductRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Height           int    `json:"height"`
	Width            int    `json:"width"`
	Thickness        int    `json:"thickness"`
	Active           bool   `json:"active"`
}

type ProductResponse struct {
	UUID             string `json:"uuid"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Active           bool   `json:"active"`
	CreatedAt        string `json:"created_at"`
}
