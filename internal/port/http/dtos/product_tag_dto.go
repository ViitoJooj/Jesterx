package dtos

type CreateProductTagRequest struct {
	ProductUUID string `json:"product_uuid"`
	Label       string `json:"label"`
}

type ProductTagResponse struct {
	UUID        string `json:"uuid"`
	ProductUUID string `json:"product_uuid"`
	Label       string `json:"label"`
	CreatedAt   string `json:"created_at"`
}
