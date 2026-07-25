package domain

import (
	"time"

	"github.com/google/uuid"
)

type Terms struct {
	UUID        uuid.UUID
	Name        string
	Description string
	UpdatedAt   *time.Time
	CreatedAt   time.Time
}
