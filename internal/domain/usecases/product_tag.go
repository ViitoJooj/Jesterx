package usecases

import (
	"database/sql"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
	"github.com/ViitoJooj/Jesterx/internal/domain/repositories/contracts"
)

type CreateProductTagUseCase struct {
	db   *sql.DB
	repo contracts.ProductTagContract
}

func NewCreateProductTagUseCase(db *sql.DB, repo contracts.ProductTagContract) *CreateProductTagUseCase {
	return &CreateProductTagUseCase{db: db, repo: repo}
}

func (u *CreateProductTagUseCase) Create(productUUID string, label string) (*domain.ProductsTags, error) {
	tag, err := domain.NewProductTag(productUUID, label, u.db)
	if err != nil {
		return nil, err
	}

	createdTag, err := u.repo.CreateProductTag(tag)
	if err != nil {
		return nil, err
	}

	return createdTag, nil
}
