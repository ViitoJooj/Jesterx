package domain

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	Id         int
	TenantID   uuid.UUID
	CreatedBy  uuid.UUID
	Name       string
	ContentMd  string
	Price      int
	Currency   string
	IsActive   bool
	Updated_at time.Time
	Created_at time.Time
}
