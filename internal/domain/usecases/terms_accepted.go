package usecases

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
	"github.com/ViitoJooj/Jesterx/internal/domain/repositories/contracts"
	"github.com/google/uuid"
)

type CreateTermsAcceptedUseCase struct {
	db                *sql.DB
	termsAcceptedRepo contracts.TermsAcceptedContract
}

func NewCreateTermsAcceptedUseCase(db *sql.DB, termsAcceptedRepo contracts.TermsAcceptedContract) *CreateTermsAcceptedUseCase {
	return &CreateTermsAcceptedUseCase{
		db:                db,
		termsAcceptedRepo: termsAcceptedRepo,
	}
}

func (u *CreateTermsAcceptedUseCase) Create(input *domain.TermsAcceptedBy, website uuid.UUID) (*domain.TermsAcceptedBy, error) {

	termsAccepted, err := domain.NewTermsAcceptedBy(
		website.String(),
		input.OwnerUUID.String(),
		string(input.OwnerType),
		u.db,
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
