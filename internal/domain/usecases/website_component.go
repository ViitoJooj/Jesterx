package usecases

import (
	"encoding/json"

	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/google/uuid"
)

type CreateWebsiteComponentUseCase struct {
	repository contracts.WebsiteComponentContract
}

func NewCreateWebsiteComponentUseCase(repository contracts.WebsiteComponentContract) *CreateWebsiteComponentUseCase {
	return &CreateWebsiteComponentUseCase{repository: repository}
}

func (u *CreateWebsiteComponentUseCase) Create(websiteUUID string, logoURL string, tittle string, description string, path string, content json.RawMessage, visits int, tenantWebsiteUUID uuid.UUID) (*domain.ComponentWebsites, error) {
	component, err := domain.NewComponentWebsite(websiteUUID, logoURL, tittle, description, path, content, visits)
	if err != nil {
		return nil, err
	}
	return u.repository.CreateWebsiteComponent(component)
}
