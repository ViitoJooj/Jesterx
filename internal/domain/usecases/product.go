package usecases

import (

	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
)

type CreateProductUseCase struct {
	repository contracts.ProductContract
}

func NewCreateProductUseCase(repository contracts.ProductContract) *CreateProductUseCase {
	return &CreateProductUseCase{repository: repository}
}

func (u *CreateProductUseCase) Create(name string, description string, shortDescription string, height int, width int, thickness int, active bool) (*domain.Products, error) {
	product, err := domain.NewProduct(name, description, shortDescription, height, width, thickness, active)
	if err != nil {
		return nil, err
	}

	createdProduct, err := u.repository.CreateProduct(product)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}
