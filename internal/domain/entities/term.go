package domain

import (
	"errors"
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

func NewTerms(name string, description string) (*Terms, error) {

	if name == "" {
		return nil, errors.New("Name cannot be null.")
	}

	if description == "" {
		return nil, errors.New("Description cannot be null.")
	}

	return &Terms{
		UUID:        uuid.Nil,
		Name:        name,
		Description: description,
	}, nil
}
