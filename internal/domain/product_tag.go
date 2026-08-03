package domain

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ViitoJooj/go-sdk/validate"
	"github.com/google/uuid"
)

type ProductsTags struct {
	UUID        uuid.UUID
	ProductUUID uuid.UUID
	Label       string
	UpdatedAt   *time.Time
	CreatedAt   time.Time
}

func NewProductTag(productUUID string, label string, db *sql.DB) (*ProductsTags, error) {

	if err := validate.UUIDv7(productUUID, "products", db); err != nil {
		return nil, err
	}

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
