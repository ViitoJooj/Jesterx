package domain

import (
	"errors"
	"time"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities/enums"
	"github.com/google/uuid"
)

type TermsAcceptedBy struct {
	UUID         uuid.UUID
	WebSiteUUID  uuid.UUID
	OwnerUUID    uuid.UUID
	OwnerType    enums.OwnerType
	AcceptedWhen time.Time
}

func NewTermsAcceptedBy(websiteUUID string, ownerUUID string, ownerType string) (*TermsAcceptedBy, error) {
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
