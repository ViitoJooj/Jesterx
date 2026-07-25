package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID        uuid.UUID
	WebSiteUUID uuid.UUID
	ImageURL    string
	Name        string
	Email       string
	Role        string
	Password    string
	CPF         *string
	UpdatedAt   *time.Time
	CreatedAt   time.Time
}
