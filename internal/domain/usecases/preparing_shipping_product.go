package usecases

import (

	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
)

type CreatePreparingShippingProductUseCase struct {
	repository contracts.PreparingShippingProductContract
}

func NewCreatePreparingShippingProductUseCase(repository contracts.PreparingShippingProductContract) *CreatePreparingShippingProductUseCase {
	return &CreatePreparingShippingProductUseCase{repository: repository}
}

func (u *CreatePreparingShippingProductUseCase) Create(productUUID string, addressUUID string) (*domain.PreparingShippingProducts, error) {
	psp, err := domain.NewPreparingShippingProduct(productUUID, addressUUID)
	if err != nil {
		return nil, err
	}

	createdPSP, err := u.repository.CreatePreparingShippingProduct(psp)
	if err != nil {
		return nil, err
	}

	return createdPSP, nil
}
