package dtos

type CreatePreparingShippingProductRequest struct {
	ProductUUID string `json:"product_uuid"`
	AddressUUID string `json:"address_uuid"`
}

type PreparingShippingProductResponse struct {
	UUID        string `json:"uuid"`
	ProductUUID string `json:"product_uuid"`
	AddressUUID string `json:"address_uuid"`
}
