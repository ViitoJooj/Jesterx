package usecases

import (
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
)

type CreateTermsUseCase struct {
	termsRepo contracts.TermsContract
}

func NewCreateTermsUseCase(termsRepo contracts.TermsContract) *CreateTermsUseCase {
	return &CreateTermsUseCase{
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

func (u *CreateTermsUseCase) GetByUUID(uuidStr string) (*domain.Terms, error) {
	return u.termsRepo.FindTermsByUUID(uuidStr)
}

func (u *CreateTermsUseCase) GetAll() ([]*domain.Terms, error) {
	return u.termsRepo.GetTerms()
}
