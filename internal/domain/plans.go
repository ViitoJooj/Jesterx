package domain

import (
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/pkg/validators"
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

func NewPlan(tenantID, createdBy uuid.UUID, name, contentMd, currency string, price int) (*Plan, error) {
	if err := validators.NewPlan(name, currency, price); err != nil {
		return nil, err
	}

	return &Plan{
		TenantID:   tenantID,
		CreatedBy:  createdBy,
		Name:       strings.TrimSpace(name),
		ContentMd:  contentMd,
		Price:      price,
		Currency:   strings.ToUpper(strings.TrimSpace(currency)),
		IsActive:   true,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}, nil
}
