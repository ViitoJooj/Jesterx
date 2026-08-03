package usecases

import (
	"database/sql"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
	"github.com/ViitoJooj/Jesterx/internal/domain/repositories/contracts"
)

type CreatePreparingShippingProductUseCase struct {
	db   *sql.DB
	repo contracts.PreparingShippingProductContract
}

func NewCreatePreparingShippingProductUseCase(db *sql.DB, repo contracts.PreparingShippingProductContract) *CreatePreparingShippingProductUseCase {
	return &CreatePreparingShippingProductUseCase{db: db, repo: repo}
}

func (u *CreatePreparingShippingProductUseCase) Create(productUUID string, addressUUID string) (*domain.PreparingShippingProducts, error) {
	psp, err := domain.NewPreparingShippingProduct(productUUID, addressUUID, u.db)
	if err != nil {
		return nil, err
	}

	createdPSP, err := u.repo.CreatePreparingShippingProduct(psp)
	if err != nil {
		return nil, err
	}

	return createdPSP, nil
}
