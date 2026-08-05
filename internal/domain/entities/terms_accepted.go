package domain

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities/enums"
	"github.com/ViitoJooj/go-sdk/validate"
	"github.com/google/uuid"
)

type TermsAcceptedBy struct {
	UUID         uuid.UUID
	WebSiteUUID  uuid.UUID
	OwnerUUID    uuid.UUID
	OwnerType    enums.OwnerType
	AcceptedWhen time.Time
}

func NewTermsAcceptedBy(websiteUUID string, ownerUUID string, ownerType string, db *sql.DB) (*TermsAcceptedBy, error) {

	if err := validate.UUIDv7(websiteUUID, "terms_accepted", db); err != nil {
		return nil, err
	}

	if err := validate.UUIDv7(ownerUUID, "users", db); err != nil {
		return nil, err
	}

	otype := enums.OwnerType(ownerType)
	if otype != enums.UserOwner && otype != enums.OrganizationOwner {
		return nil, errors.New("OwnerType must be 'User' or 'Organization'.")
	}

	websiteUUIDParsed, err := uuid.Parse(websiteUUID)
	if err != nil {
		return nil, err
	}

	ownerUUIDParsed, err := uuid.Parse(ownerUUID)
	if err != nil {
		return nil, err
	}

	return &TermsAcceptedBy{
		UUID:         uuid.Nil,
		WebSiteUUID:  websiteUUIDParsed,
		OwnerUUID:    ownerUUIDParsed,
		OwnerType:    otype,
		AcceptedWhen: time.Now(),
	}, nil
}
