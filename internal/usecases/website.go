package usecases

import (
	"database/sql"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repositories/contracts"
)

type CreateWebsiteUseCase struct {
	db   *sql.DB
	repo contracts.WebsiteContract
}

func NewCreateWebsiteUseCase(db *sql.DB, repo contracts.WebsiteContract) *CreateWebsiteUseCase {
	return &CreateWebsiteUseCase{db: db, repo: repo}
}

func (u *CreateWebsiteUseCase) Create(ownerUUID string, ownerType string, label string, url string, writeIn string, description string) (*domain.Website, error) {
	website, err := domain.NewWebsite(ownerUUID, ownerType, label, url, writeIn, description, u.db)
	if err != nil {
		return nil, err
	}
	return u.repo.CreateWebsite(website)
}
