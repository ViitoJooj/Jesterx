package domain

import (
	"time"

	"github.com/google/uuid"
)

type Rbac struct {
	UUID        uuid.UUID
	WebSiteUUID uuid.UUID
	Label       string
	CanRead     bool
	CanWrite    bool
	CanUpdate   bool
	CanUpgrade  bool
	CanDelete   bool
	UpdatedAt   *time.Time
	CreatedAt   time.Time
}
