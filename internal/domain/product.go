package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Products struct {
	UUID             uuid.UUID
	Name             string
	Description      string
	ShortDescription string
	height           int
	width            int
	thickness        int
	Active           bool
	UpdatedAt        *time.Time
	CreatedAt        time.Time
}

func NewProduct(name string, description string, shortDescription string, height int, width int, thickness int, active bool) (*Products, error) {

	if name == "" {
		return nil, errors.New("Name cannot be null.")
	}

	if height < 0 {
		return nil, errors.New("Height cannot be negative.")
	}

	if width < 0 {
		return nil, errors.New("Width cannot be negative.")
	}

	if thickness < 0 {
		return nil, errors.New("Thickness cannot be negative.")
	}

	return &Products{
		UUID:             uuid.Nil,
		Name:             name,
		Description:      description,
		ShortDescription: shortDescription,
		height:           height,
		width:            width,
		thickness:        thickness,
		Active:           active,
	}, nil
}
