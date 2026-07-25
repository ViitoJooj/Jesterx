package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationBR struct {
	UUID        uuid.UUID
	WebSiteUUID uuid.UUID
	OwnerUUID   uuid.UUID
	ImageURL    string
	Name        string
	TradeName   string
	CNPJ        string
	UpdatedAt   *time.Time
	CreatedAt   time.Time
}
