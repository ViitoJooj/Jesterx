package usecases

import (
	"database/sql"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repositories/contracts"
)

type CreateStorageProductUseCase struct {
	db   *sql.DB
	repo contracts.StorageProductContract
}

func NewCreateStorageProductUseCase(db *sql.DB, repo contracts.StorageProductContract) *CreateStorageProductUseCase {
	return &CreateStorageProductUseCase{db: db, repo: repo}
}

func (u *CreateStorageProductUseCase) Create(productUUID string) (*domain.StorageProducts, error) {
	sp, err := domain.NewStorageProduct(productUUID, u.db)
	if err != nil {
		return nil, err
	}

	createdSP, err := u.repo.CreateStorageProduct(sp)
	if err != nil {
		return nil, err
	}

	return createdSP, nil
}
