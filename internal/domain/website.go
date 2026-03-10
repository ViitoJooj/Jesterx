package domain

import (
	"time"

	"github.com/google/uuid"
)

type WebSite struct {
	Id                string
	Type              string
	Image             []byte
	Name              string
	Short_description string
	Description       string
	Creator_id        string
	Banned            bool
	MatureContent     bool
	RatingAvg         float64
	RatingCount       int
	Updated_at        time.Time
	Created_at        time.Time
}

func NewWebSite(Type string, Image []byte, Name string, Short_description string, Description string, Creator_id string) *WebSite {
	id, _ := uuid.NewV7()

	return &WebSite{
		Id:                id.String(),
		Type:              Type,
		Image:             Image,
		Name:              Name,
		Short_description: Short_description,
		Description:       Description,
		Creator_id:        Creator_id,
		Banned:            false,
		Updated_at:        time.Now(),
		Created_at:        time.Now(),
	}
}

type WebSiteRoute struct {
	Id           string
	WebsiteId    string
	Path         string
	Title        string
	RequiresAuth bool
	Position     int
	Updated_at   time.Time
	Created_at   time.Time
}

func NewWebSiteRoute(websiteID string, path string, title string, requiresAuth bool, position int) *WebSiteRoute {
	id, _ := uuid.NewV7()
	now := time.Now()

	return &WebSiteRoute{
		Id:           id.String(),
		WebsiteId:    websiteID,
		Path:         path,
		Title:        title,
		RequiresAuth: requiresAuth,
		Position:     position,
		Updated_at:   now,
		Created_at:   now,
	}
}

type WebSiteVersion struct {
	Id           string
	WebsiteId    string
	Version      int
	SourceType   string
	Source       string
	CompiledHTML string
	ScanStatus   string
	ScanScore    int
	ScanFindings string
	Published    bool
	PublishedAt  *time.Time
	CreatedBy    string
	Updated_at   time.Time
	Created_at   time.Time
}

func NewWebSiteVersion(websiteID string, version int, sourceType string, source string, createdBy string) *WebSiteVersion {
	id, _ := uuid.NewV7()
	now := time.Now()
	return &WebSiteVersion{
		Id:         id.String(),
		WebsiteId:  websiteID,
		Version:    version,
		SourceType: sourceType,
		Source:     source,
		CreatedBy:  createdBy,
		Updated_at: now,
		Created_at: now,
	}
}
