package repository

import "github.com/ViitoJooj/Jesterx/internal/domain"

type WebsiteRepository interface {
	SaveWebSite(website domain.WebSite) error
	FindWebSiteByID(id string) (*domain.WebSite, error)
	FindWebSiteByUserID(id string) (*domain.WebSite, error)
	FindWebSiteByName(name string) (*domain.WebSite, error)
	UpdateWebSiteByID(id string, website domain.WebSite) error
	DeleteWebSiteByID(id string) error
}
