package repository

import (
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

type WebsiteRepository interface {
	SaveWebSite(website domain.WebSite) error
	FindWebSiteByID(id string) (*domain.WebSite, error)
	ListWebSitesByUserID(id string) ([]domain.WebSite, error)
	FindWebSiteByName(name string) (*domain.WebSite, error)
	UpdateWebSiteByID(id string, website domain.WebSite) error
	DeleteWebSiteByID(id string) error
	CountWebSitesByUserID(userID string) (int, error)
	ReplaceRoutesByWebsiteID(websiteID string, routes []domain.WebSiteRoute) error
	ListRoutesByWebsiteID(websiteID string) ([]domain.WebSiteRoute, error)
	FindRouteByWebsiteIDAndPath(websiteID string, path string) (*domain.WebSiteRoute, error)
	SaveVersion(version domain.WebSiteVersion) error
	DeleteVersionsByWebsiteID(websiteID string) error
	FindLatestVersionByWebsiteID(websiteID string) (*domain.WebSiteVersion, error)
	ListVersionsByWebsiteID(websiteID string) ([]domain.WebSiteVersion, error)
	FindVersionByWebsiteID(websiteID string, version int) (*domain.WebSiteVersion, error)
	FindPublishedVersionByWebsiteID(websiteID string) (*domain.WebSiteVersion, error)
	UpdateVersionPublishState(websiteID string, version int, published bool, publishedAt *time.Time) error
}
