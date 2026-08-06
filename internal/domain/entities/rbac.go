package domain

import (
	"errors"
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

func NewRbac(websiteUUID string, label string, canRead bool, canWrite bool, canUpdate bool, canUpgrade bool, canDelete bool) (*Rbac, error) {
	if label == "" {
		return nil, errors.New("Label cannot be null.")
	}

	websiteUUIDParsed, err := uuid.Parse(websiteUUID)
	if err != nil {
		return nil, err
	}

	return &Rbac{
		UUID:        uuid.Nil,
		WebSiteUUID: websiteUUIDParsed,
		Label:       label,
		CanRead:     canRead,
		CanWrite:    canWrite,
		CanUpdate:   canUpdate,
		CanUpgrade:  canUpgrade,
		CanDelete:   canDelete,
	}, nil
}
