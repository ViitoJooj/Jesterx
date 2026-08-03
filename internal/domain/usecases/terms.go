package usecases

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
	"github.com/ViitoJooj/Jesterx/internal/domain/repositories/contracts"
)

type CreateTermsUseCase struct {
	db        *sql.DB
	termsRepo contracts.TermsContract
}

func NewCreateTermsUseCase(db *sql.DB, termsRepo contracts.TermsContract) *CreateTermsUseCase {
	return &CreateTermsUseCase{
		db:        db,
		termsRepo: termsRepo,
	}
}

func (u *CreateTermsUseCase) Create(input *domain.Terms) (*domain.Terms, error) {

	terms, err := domain.NewTerms(input.Name, input.Description)
	if err != nil {
		return nil, err
	}

	existing, err := u.termsRepo.FindTermsByName(input.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("terms already exists for this name")
	}

	createdTerms, err := u.termsRepo.CreateTerms(terms)
	if err != nil {
		return nil, err
	}

	return createdTerms, nil
}
