package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	Id        int64
	TenantID  uuid.UUID
	CreatedBy uuid.UUID
	Type      string
	IsDefault bool
	UpdatedAt time.Time
	CreatedAt time.Time
}
