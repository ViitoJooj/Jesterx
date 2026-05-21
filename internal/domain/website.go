package domain

import (
	"time"

	"github.com/google/uuid"
)

type WebSite struct {
	Id                string
	Type              string
	ImageUrl          *string
	Name              string
	Short_description string
	Description       string
	Creator_id        string
	CompanyId         *string
	Banned            bool
	MatureContent     bool
	RatingAvg         float64
	RatingCount       int
	Updated_at        time.Time
	Created_at        time.Time
}

func NewWebSite(siteType string, imageUrl *string, name, shortDesc, description, creatorId string) *WebSite {
	id, _ := uuid.NewV7()
	now := time.Now()
	return &WebSite{
		Id:                id.String(),
		Type:              siteType,
		ImageUrl:          imageUrl,
		Name:              name,
		Short_description: shortDesc,
		Description:       description,
		Creator_id:        creatorId,
		Banned:            false,
		Updated_at:        now,
		Created_at:        now,
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

func NewWebSiteRoute(websiteID, path, title string, requiresAuth bool, position int) *WebSiteRoute {
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

func NewWebSiteVersion(websiteID string, version int, sourceType, source, createdBy string) *WebSiteVersion {
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
