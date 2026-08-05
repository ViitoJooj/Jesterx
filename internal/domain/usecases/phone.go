package usecases

import (
	"database/sql"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/google/uuid"
)

type CreatePhoneUseCase struct {
	db        *sql.DB
	phoneRepo contracts.PhoneContract
}

func NewCreatePhoneUseCase(db *sql.DB, phoneRepo contracts.PhoneContract) *CreatePhoneUseCase {
	return &CreatePhoneUseCase{
		db:        db,
		phoneRepo: phoneRepo,
	}
}

func (u *CreatePhoneUseCase) Create(input *domain.Phone, website uuid.UUID) (*domain.Phone, error) {

	phone, err := domain.NewPhone(
		website.String(),
		input.OwnerUUID.String(),
		string(input.OwnerType),
		input.Label,
		input.Number,
		input.IsDefault,
		u.db,
	)
	if err != nil {
		return nil, err
	}

	createdPhone, err := u.phoneRepo.CreatePhone(phone)
	if err != nil {
		return nil, err
	}

	return createdPhone, nil
}
