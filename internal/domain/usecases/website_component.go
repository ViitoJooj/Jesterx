package usecases

import (
	"database/sql"
	"encoding/json"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/google/uuid"
)

type CreateWebsiteComponentUseCase struct {
	db   *sql.DB
	repo contracts.WebsiteComponentContract
}

func NewCreateWebsiteComponentUseCase(db *sql.DB, repo contracts.WebsiteComponentContract) *CreateWebsiteComponentUseCase {
	return &CreateWebsiteComponentUseCase{db: db, repo: repo}
}

func (u *CreateWebsiteComponentUseCase) Create(websiteUUID string, logoURL string, tittle string, description string, path string, content json.RawMessage, visits int, tenantWebsiteUUID uuid.UUID) (*domain.ComponentWebsites, error) {
	component, err := domain.NewComponentWebsite(websiteUUID, logoURL, tittle, description, path, content, visits, u.db)
	if err != nil {
		return nil, err
	}
	return u.repo.CreateWebsiteComponent(component)
}
