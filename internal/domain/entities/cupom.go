package domain

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities/enums"
	"github.com/ViitoJooj/go-sdk/validate"
	"github.com/google/uuid"
)

type Cupons struct {
	UUID        uuid.UUID
	TagUUID     uuid.UUID
	Label       string
	Description string
	Value       string
	ValueType   enums.CupomValueType
}

func NewCupom(tagUUID string, label string, description string, value string, valueType string, db *sql.DB) (*Cupons, error) {

	if err := validate.UUIDv7(tagUUID, "products_tags", db); err != nil {
		return nil, err
	}

	if label == "" {
		return nil, errors.New("Label cannot be null.")
	}

	if value == "" {
		return nil, errors.New("Value cannot be null.")
	}

	vtype := enums.CupomValueType(valueType)
	if vtype != enums.CupomPercentage && vtype != enums.CupomValue {
		return nil, errors.New("ValueType must be 'percentage' or 'Value'.")
	}

	tagUUIDParsed, err := uuid.Parse(tagUUID)
	if err != nil {
		return nil, err
	}

	return &Cupons{
		UUID:        uuid.Nil,
		TagUUID:     tagUUIDParsed,
		Label:       label,
		Description: description,
		Value:       value,
		ValueType:   vtype,
	}, nil
}
