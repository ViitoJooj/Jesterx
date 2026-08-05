package domain

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities/enums"
	"github.com/ViitoJooj/go-sdk/validate"
	"github.com/google/uuid"
)

type Phone struct {
	UUID        uuid.UUID
	WebSiteUUID uuid.UUID
	OwnerUUID   uuid.UUID
	OwnerType   enums.OwnerType
	Label       string
	Number      int
	IsDefault   bool
	UpdatedAt   *time.Time
	CreatedAt   time.Time
}

func NewPhone(websiteUUID string, ownerUUID string, ownerType string, label string, number int, isDefault bool, db *sql.DB) (*Phone, error) {

	if err := validate.UUIDv7(websiteUUID, "phones", db); err != nil {
		return nil, err
	}

	if err := validate.UUIDv7(ownerUUID, "users", db); err != nil {
		return nil, err
	}

	otype := enums.OwnerType(ownerType)
	if otype != enums.UserOwner && otype != enums.OrganizationOwner {
		return nil, errors.New("OwnerType must be 'User' or 'Organization'.")
	}

	if label == "" {
		return nil, errors.New("Label cannot be null.")
	}

	websiteUUIDParsed, err := uuid.Parse(websiteUUID)
	if err != nil {
		return nil, err
	}

	ownerUUIDParsed, err := uuid.Parse(ownerUUID)
	if err != nil {
		return nil, err
	}

	return &Phone{
		UUID:        uuid.Nil,
		WebSiteUUID: websiteUUIDParsed,
		OwnerUUID:   ownerUUIDParsed,
		OwnerType:   otype,
		Label:       label,
		Number:      number,
		IsDefault:   isDefault,
	}, nil
}
