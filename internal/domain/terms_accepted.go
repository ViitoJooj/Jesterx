package domain

import (
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain/enums"
	"github.com/google/uuid"
)

type TermsAcceptedBy struct {
	UUID         uuid.UUID
	WebSiteUUID  uuid.UUID
	OwnerUUID    uuid.UUID
	OwnerType    enums.OwnerType
	AcceptedWhen time.Time
}
