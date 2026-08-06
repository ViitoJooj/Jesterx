package usecases

import (
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/google/uuid"
)

type CreateTermsAcceptedUseCase struct {
	termsAcceptedRepo contracts.TermsAcceptedContract
}

func NewCreateTermsAcceptedUseCase(termsAcceptedRepo contracts.TermsAcceptedContract) *CreateTermsAcceptedUseCase {
	return &CreateTermsAcceptedUseCase{
		termsAcceptedRepo: termsAcceptedRepo,
	}
}

func (u *CreateTermsAcceptedUseCase) Create(input *domain.TermsAcceptedBy, website uuid.UUID) (*domain.TermsAcceptedBy, error) {

	termsAccepted, err := domain.NewTermsAcceptedBy(
		website.String(),
		input.OwnerUUID.String(),
		string(input.OwnerType),
	)
	if err != nil {
		return nil, err
	}

	existing, err := u.termsAcceptedRepo.FindTermsAcceptedByOwner(input.OwnerUUID.String(), website.String())
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("terms already accepted by this owner")
	}

	createdTermsAccepted, err := u.termsAcceptedRepo.CreateTermsAccepted(termsAccepted)
	if err != nil {
		return nil, err
	}

	return createdTermsAccepted, nil
}

func (u *CreateTermsAcceptedUseCase) GetByUUID(uuidStr string) (*domain.TermsAcceptedBy, error) {
	return u.termsAcceptedRepo.FindTermsAcceptedByUUID(uuidStr)
}

func (u *CreateTermsAcceptedUseCase) GetAll(websiteUUIDStr string) ([]*domain.TermsAcceptedBy, error) {
	return u.termsAcceptedRepo.GetTermsAcceptedFromWebsite(websiteUUIDStr)
}
