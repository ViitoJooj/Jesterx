package validators

import (
	"errors"
	"strings"
)

var validPermissions = map[string]struct{}{
	"read":   {},
	"write":  {},
	"delete": {},
	"manage": {},
}

func NewRole(name string, permissions []string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("role name is required.")
	}

	for _, p := range permissions {
		if _, ok := validPermissions[p]; !ok {
			return errors.New("invalid permission: " + p)
		}
	}

	return nil
}
