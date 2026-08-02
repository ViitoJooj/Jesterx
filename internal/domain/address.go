package domain

import (
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain/enums"
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

func NewAddressesBR(
	websiteUUID,
	ownerUUID uuid.UUID,
	ownerType enums.OwnerType,
	label string,
	addressline1 string,
	addressline2 string,
	neighborhood string,
	city string,
	state string,
	statecCode string,
	postalCode string,
	referencePoint string,
	DeliveryNotes string) (*AddressBR, error) {

	output := AddressBR{
		UUID:           uuid.Nil,
		WebSiteUUID:    websiteUUID,
		OwnerUUID:      ownerUUID,
		OwnerType:      ownerType,
		Label:          label,
		AddressLine1:   addressline1,
		AddressLine2:   addressline2,
		City:           city,
		State:          state,
		StateCode:      statecCode,
		PostalCode:     postalCode,
		ReferencePoint: referencePoint,
		DeliveryNotes:  DeliveryNotes,
		IsDefault:      false,
		UpdatedAt:      nil,
		CreatedAt:      time.Now(),
	}

	return &output, nil
}
