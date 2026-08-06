package usecases

import (

	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
)

type CreateProductShippedUseCase struct {
	repository contracts.ProductShippedContract
}

func NewCreateProductShippedUseCase(repository contracts.ProductShippedContract) *CreateProductShippedUseCase {
	return &CreateProductShippedUseCase{repository: repository}
}

func (u *CreateProductShippedUseCase) Create(productUUID string, addressUUID string, status string) (*domain.ProductShipped, error) {
	productShipped, err := domain.NewProductShipped(productUUID, addressUUID, status)
	if err != nil {
		return nil, err
	}
	return u.repository.CreateProductShipped(productShipped)
}
