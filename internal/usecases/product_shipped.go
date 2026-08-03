package usecases

import (
	"database/sql"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repositories/contracts"
)

type CreateProductShippedUseCase struct {
	db   *sql.DB
	repo contracts.ProductShippedContract
}

func NewCreateProductShippedUseCase(db *sql.DB, repo contracts.ProductShippedContract) *CreateProductShippedUseCase {
	return &CreateProductShippedUseCase{db: db, repo: repo}
}

func (u *CreateProductShippedUseCase) Create(productUUID string, addressUUID string, status string) (*domain.ProductShipped, error) {
	productShipped, err := domain.NewProductShipped(productUUID, addressUUID, status, u.db)
	if err != nil {
		return nil, err
	}
	return u.repo.CreateProductShipped(productShipped)
}
