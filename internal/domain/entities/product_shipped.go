package domain

import (
	"github.com/google/uuid"
)

type ProductShipped struct {
	UUID        uuid.UUID
	ProductUUID uuid.UUID
	AddressUUID uuid.UUID
	Status      string
}

func NewProductShipped(productUUID string, addressUUID string, status string) (*ProductShipped, error) {
	productUUIDParsed, err := uuid.Parse(productUUID)
	if err != nil {
		return nil, err
	}

	addressUUIDParsed, err := uuid.Parse(addressUUID)
	if err != nil {
		return nil, err
	}

	return &ProductShipped{
		UUID:        uuid.Nil,
		ProductUUID: productUUIDParsed,
		AddressUUID: addressUUIDParsed,
		Status:      status,
	}, nil
}
