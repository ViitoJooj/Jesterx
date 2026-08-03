package dtos

type CreateStorageProductRequest struct {
	ProductUUID string `json:"product_uuid"`
}

type StorageProductResponse struct {
	UUID        string `json:"uuid"`
	ProductUUID string `json:"product_uuid"`
}
