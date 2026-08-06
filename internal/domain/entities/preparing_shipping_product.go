package domain

import (
	"github.com/google/uuid"
)

type PreparingShippingProducts struct {
	UUID        uuid.UUID
	ProductUUID uuid.UUID
	AddressUUID uuid.UUID
}

func NewPreparingShippingProduct(productUUID string, addressUUID string) (*PreparingShippingProducts, error) {
	productUUIDParsed, err := uuid.Parse(productUUID)
	if err != nil {
		return nil, err
	}

	addressUUIDParsed, err := uuid.Parse(addressUUID)
	if err != nil {
		return nil, err
	}

	return &PreparingShippingProducts{
		UUID:        uuid.Nil,
		ProductUUID: productUUIDParsed,
		AddressUUID: addressUUIDParsed,
	}, nil
}
