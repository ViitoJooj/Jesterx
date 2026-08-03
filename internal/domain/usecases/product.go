package usecases

import (
	"database/sql"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
	"github.com/ViitoJooj/Jesterx/internal/domain/repositories/contracts"
)

type CreateProductUseCase struct {
	db   *sql.DB
	repo contracts.ProductContract
}

func NewCreateProductUseCase(db *sql.DB, repo contracts.ProductContract) *CreateProductUseCase {
	return &CreateProductUseCase{db: db, repo: repo}
}

func (u *CreateProductUseCase) Create(name string, description string, shortDescription string, height int, width int, thickness int, active bool) (*domain.Products, error) {
	product, err := domain.NewProduct(name, description, shortDescription, height, width, thickness, active)
	if err != nil {
		return nil, err
	}

	createdProduct, err := u.repo.CreateProduct(product)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}
