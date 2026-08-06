package domain

import (
	"github.com/google/uuid"
)

type StorageProducts struct {
	UUID        uuid.UUID
	ProductUUID uuid.UUID
}

func NewStorageProduct(productUUID string) (*StorageProducts, error) {
	productUUIDParsed, err := uuid.Parse(productUUID)
	if err != nil {
		return nil, err
	}

	return &StorageProducts{
		UUID:        uuid.Nil,
		ProductUUID: productUUIDParsed,
	}, nil
}
