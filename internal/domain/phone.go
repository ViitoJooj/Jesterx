package domain

import (
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain/enums"
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
