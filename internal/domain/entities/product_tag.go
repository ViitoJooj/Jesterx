package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ProductsTags struct {
	UUID        uuid.UUID
	ProductUUID uuid.UUID
	Label       string
	UpdatedAt   *time.Time
	CreatedAt   time.Time
}

func NewProductTag(productUUID string, label string) (*ProductsTags, error) {
	if label == "" {
		return nil, errors.New("Label cannot be null.")
	}

	productUUIDParsed, err := uuid.Parse(productUUID)
	if err != nil {
		return nil, err
	}

	return &ProductsTags{
		UUID:        uuid.Nil,
		ProductUUID: productUUIDParsed,
		Label:       label,
	}, nil
}
