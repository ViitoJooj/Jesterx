package domain

import (
	"time"

	"github.com/ViitoJooj/Jesterx/pkg/validators"
	"github.com/google/uuid"
)

type Permission string

const (
	PermissionRead   Permission = "read"
	PermissionWrite  Permission = "write"
	PermissionDelete Permission = "delete"
	PermissionManage Permission = "manage"
)

type Role struct {
	Id          int64
	TenantID    uuid.UUID
	CreatedBy   uuid.UUID
	Name        string
	Permissions []Permission
	IsDefault   bool
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

func NewRole(tenantID, createdBy uuid.UUID, name string, permissions []Permission, isDefault bool) (*Role, error) {
	perms := make([]string, len(permissions))
	for i, p := range permissions {
		perms[i] = string(p)
	}

	if err := validators.NewRole(name, perms); err != nil {
		return nil, err
	}

	return &Role{
		TenantID:    tenantID,
		CreatedBy:   createdBy,
		Name:        name,
		Permissions: permissions,
		IsDefault:   isDefault,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
