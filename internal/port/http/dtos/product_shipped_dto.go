package dtos

type CreateProductShippedRequest struct {
	ProductUUID string `json:"product_uuid"`
	AddressUUID string `json:"address_uuid"`
	Status      string `json:"status"`
}

type ProductShippedResponse struct {
	UUID        string `json:"uuid"`
	ProductUUID string `json:"product_uuid"`
	AddressUUID string `json:"address_uuid"`
	Status      string `json:"status"`
}
