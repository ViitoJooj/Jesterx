package domain

import (
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain/enums"
	"github.com/google/uuid"
)

type UserAddressBR struct {
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
