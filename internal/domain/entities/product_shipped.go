package domain

import (
	"database/sql"

	"github.com/ViitoJooj/go-sdk/validate"
	"github.com/google/uuid"
)

type ProductShipped struct {
	UUID        uuid.UUID
	ProductUUID uuid.UUID
	AddressUUID uuid.UUID
	Status      string
}

func NewProductShipped(productUUID string, addressUUID string, status string, db *sql.DB) (*ProductShipped, error) {

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

	return &ProductShipped{
		UUID:        uuid.Nil,
		ProductUUID: productUUIDParsed,
		AddressUUID: addressUUIDParsed,
		Status:      status,
	}, nil
}
