package domain

import (
	"errors"
	"time"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities/enums"
	"github.com/google/uuid"
)

type AddressBR struct {
	UUID           uuid.UUID
	WebSiteUUID    uuid.UUID
	OwnerUUID      uuid.UUID
	OwnerType      enums.OwnerType
	Label          string
	AddressLine1   string
	AddressLine2   string
	Neighborhood   string
	City           string
	State          string
	StateCode      string
	PostalCode     string
	ReferencePoint string
	DeliveryNotes  string
	IsDefault      bool
	UpdatedAt      *time.Time
	CreatedAt      time.Time
}

func NewAddress(
	websiteUUID string,
	ownerUUID string,
	ownerType string,
	label string,
	addressLine1 string,
	addressLine2 string,
	neighborhood string,
	city string,
	state string,
	stateCode string,
	postalCode string,
	referencePoint string,
	deliveryNotes string,
	isDefault bool,
) (*AddressBR, error) {

	otype := enums.OwnerType(ownerType)
	if otype != enums.UserOwner && otype != enums.OrganizationOwner {
		return nil, errors.New("OwnerType must be 'User' or 'Organization'.")
	}

	if label == "" {
		return nil, errors.New("Label cannot be null.")
	}

	if city == "" {
		return nil, errors.New("City cannot be null.")
	}

	if state == "" {
		return nil, errors.New("State cannot be null.")
	}

	if len(stateCode) != 2 {
		return nil, errors.New("StateCode must be 2 characters.")
	}

	if postalCode == "" {
		return nil, errors.New("PostalCode cannot be null.")
	}

	websiteUUIDParsed, err := uuid.Parse(websiteUUID)
	if err != nil {
		return nil, err
	}

	ownerUUIDParsed, err := uuid.Parse(ownerUUID)
	if err != nil {
		return nil, err
	}

	return &AddressBR{
		UUID:           uuid.Nil,
		WebSiteUUID:    websiteUUIDParsed,
		OwnerUUID:      ownerUUIDParsed,
		OwnerType:      otype,
		Label:          label,
		AddressLine1:   addressLine1,
		AddressLine2:   addressLine2,
		Neighborhood:   neighborhood,
		City:           city,
		State:          state,
		StateCode:      stateCode,
		PostalCode:     postalCode,
		ReferencePoint: referencePoint,
		DeliveryNotes:  deliveryNotes,
		IsDefault:      isDefault,
	}, nil
}
