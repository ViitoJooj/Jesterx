package domain

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities/enums"
	"github.com/ViitoJooj/go-sdk/validate"
	"github.com/google/uuid"
)

type Website struct {
	UUID        uuid.UUID
	OwnerUUID   uuid.UUID
	OwnerType   enums.OwnerType
	Label       string
	URL         string
	WriteIn     enums.LanguageType
	Description string
	UpdatedAt   *time.Time
	CreatedAt   time.Time
}

func NewWebsite(ownerUUID string, ownerType string, label string, url string, writeIn string, description string, db *sql.DB) (*Website, error) {

	otype := enums.OwnerType(ownerType)
	if otype != enums.UserOwner && otype != enums.OrganizationOwner {
		return nil, errors.New("OwnerType must be 'User' or 'Organization'.")
	}

	var ownerTable string
	if otype == enums.UserOwner {
		ownerTable = "users"
	} else {
		ownerTable = "organizations"
	}

	if err := validate.UUIDv7(ownerUUID, ownerTable, db); err != nil {
		return nil, err
	}

	if label == "" {
		return nil, errors.New("Label cannot be null.")
	}

	lang := enums.LanguageType(writeIn)
	if lang != enums.Component && lang != enums.React && lang != enums.Svelte {
		return nil, errors.New("WriteIn must be 'Component', 'React' or 'Svelte'.")
	}

	ownerUUIDParsed, err := uuid.Parse(ownerUUID)
	if err != nil {
		return nil, err
	}

	return &Website{
		UUID:        uuid.Nil,
		OwnerUUID:   ownerUUIDParsed,
		OwnerType:   otype,
		Label:       label,
		URL:         url,
		WriteIn:     lang,
		Description: description,
	}, nil
}
