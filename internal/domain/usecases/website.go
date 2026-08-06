package usecases

import (

	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
)

type CreateWebsiteUseCase struct {
	repository contracts.WebsiteContract
}

func NewCreateWebsiteUseCase(repository contracts.WebsiteContract) *CreateWebsiteUseCase {
	return &CreateWebsiteUseCase{repository: repository}
}

func (u *CreateWebsiteUseCase) Create(ownerUUID string, ownerType string, label string, url string, writeIn string, description string) (*domain.Website, error) {
	website, err := domain.NewWebsite(ownerUUID, ownerType, label, url, writeIn, description)
	if err != nil {
		return nil, err
	}
	return u.repository.CreateWebsite(website)
}

func (u *CreateWebsiteUseCase) GetByUUID(uuidStr string) (*domain.Website, error) {
	return u.repository.FindWebsiteByUUID(uuidStr)
}

func (u *CreateWebsiteUseCase) ListByOwner(userUUID string) ([]*domain.Website, error) {
	return u.repository.FindWebsitesByOwner(userUUID)
}

func (u *CreateWebsiteUseCase) ListAll() ([]*domain.Website, error) {
	return u.repository.GetWebsites()
}

func (u *CreateWebsiteUseCase) Update(uuidStr string) error {
	return u.repository.UpdateWebsiteByUUID(uuidStr)
}

func (u *CreateWebsiteUseCase) Delete(uuidStr string) error {
	return u.repository.DeleteWebsiteByUUID(uuidStr)
}
