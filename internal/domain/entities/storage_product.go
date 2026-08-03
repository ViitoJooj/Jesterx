package domain

import (
	"database/sql"

	"github.com/ViitoJooj/go-sdk/validate"
	"github.com/google/uuid"
)

type StorageProducts struct {
	UUID        uuid.UUID
	ProductUUID uuid.UUID
}

func NewStorageProduct(productUUID string, db *sql.DB) (*StorageProducts, error) {

	if err := validate.UUIDv7(productUUID, "products", db); err != nil {
		return nil, err
	}

	productUUIDParsed, err := uuid.Parse(productUUID)
	if err != nil {
		return nil, err
	}

	return &StorageProducts{
		UUID:        uuid.Nil,
		ProductUUID: productUUIDParsed,
	}, nil
}
