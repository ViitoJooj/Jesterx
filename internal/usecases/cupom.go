package usecases

import (
	"database/sql"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repositories/contracts"
)

type CreateCupomUseCase struct {
	db   *sql.DB
	repo contracts.CupomContract
}

func NewCreateCupomUseCase(db *sql.DB, repo contracts.CupomContract) *CreateCupomUseCase {
	return &CreateCupomUseCase{db: db, repo: repo}
}

func (u *CreateCupomUseCase) Create(tagUUID string, label string, description string, value string, valueType string) (*domain.Cupons, error) {
	cupom, err := domain.NewCupom(tagUUID, label, description, value, valueType, u.db)
	if err != nil {
		return nil, err
	}

	createdCupom, err := u.repo.CreateCupom(cupom)
	if err != nil {
		return nil, err
	}

	return createdCupom, nil
}
