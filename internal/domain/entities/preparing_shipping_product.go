package domain

import (
	"database/sql"

	"github.com/ViitoJooj/go-sdk/validate"
	"github.com/google/uuid"
)

type PreparingShippingProducts struct {
	UUID        uuid.UUID
	ProductUUID uuid.UUID
	AddressUUID uuid.UUID
}

func NewPreparingShippingProduct(productUUID string, addressUUID string, db *sql.DB) (*PreparingShippingProducts, error) {

	if err := validate.UUIDv7(productUUID, "products", db); err != nil {
		return nil, err
	}

	if err := validate.UUIDv7(addressUUID, "addresses", db); err != nil {
		return nil, err
	}

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
