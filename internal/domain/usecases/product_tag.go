package usecases

import (

	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
)

type CreateProductTagUseCase struct {
	repository contracts.ProductTagContract
}

func NewCreateProductTagUseCase(repository contracts.ProductTagContract) *CreateProductTagUseCase {
	return &CreateProductTagUseCase{repository: repository}
}

func (u *CreateProductTagUseCase) Create(productUUID string, label string) (*domain.ProductsTags, error) {
	tag, err := domain.NewProductTag(productUUID, label)
	if err != nil {
		return nil, err
	}

	createdTag, err := u.repository.CreateProductTag(tag)
	if err != nil {
		return nil, err
	}

	return createdTag, nil
}
