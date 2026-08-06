package usecases

import (

	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
)

type CreateStorageProductUseCase struct {
	repository contracts.StorageProductContract
}

func NewCreateStorageProductUseCase(repository contracts.StorageProductContract) *CreateStorageProductUseCase {
	return &CreateStorageProductUseCase{repository: repository}
}

func (u *CreateStorageProductUseCase) Create(productUUID string) (*domain.StorageProducts, error) {
	sp, err := domain.NewStorageProduct(productUUID)
	if err != nil {
		return nil, err
	}

	createdSP, err := u.repository.CreateStorageProduct(sp)
	if err != nil {
		return nil, err
	}

	return createdSP, nil
}
