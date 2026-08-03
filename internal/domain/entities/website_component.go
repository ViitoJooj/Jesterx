package domain

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/ViitoJooj/go-sdk/validate"
	"github.com/google/uuid"
)

type ComponentWebsites struct {
	UUID        uuid.UUID
	WebsiteUUID uuid.UUID
	LogoURL     string
	Tittle      string
	Description string
	Path        string
	Content     json.RawMessage
	Visists     int
	UpdatedBy   *uuid.UUID
	UpdatedAt   *time.Time
	CreatedAt   time.Time
}

func NewComponentWebsite(websiteUUID string, logoURL string, tittle string, description string, path string, content json.RawMessage, visits int, db *sql.DB) (*ComponentWebsites, error) {

	if err := validate.UUIDv7(websiteUUID, "websites", db); err != nil {
		return nil, err
	}

	if tittle == "" {
		return nil, errors.New("Tittle cannot be null.")
	}

	websiteUUIDParsed, err := uuid.Parse(websiteUUID)
	if err != nil {
		return nil, err
	}

	return &ComponentWebsites{
		UUID:        uuid.Nil,
		WebsiteUUID: websiteUUIDParsed,
		LogoURL:     logoURL,
		Tittle:      tittle,
		Description: description,
		Path:        path,
		Content:     content,
		Visists:     visits,
	}, nil
}
