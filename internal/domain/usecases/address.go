package usecases

import (
	"database/sql"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
	"github.com/ViitoJooj/Jesterx/internal/domain/repositories/contracts"
	"github.com/google/uuid"
)

type CreateAddressUseCase struct {
	db          *sql.DB
	addressRepo contracts.AddressContract
}

func NewCreateAddressUseCase(db *sql.DB, addressRepo contracts.AddressContract) *CreateAddressUseCase {
	return &CreateAddressUseCase{
		db:          db,
		addressRepo: addressRepo,
	}
}

func (u *CreateAddressUseCase) Create(input *domain.AddressBR, website uuid.UUID) (*domain.AddressBR, error) {

	address, err := domain.NewAddress(
		website.String(),
		input.OwnerUUID.String(),
		string(input.OwnerType),
		input.Label,
		input.AddressLine1,
		input.AddressLine2,
		input.Neighborhood,
		input.City,
		input.State,
		input.StateCode,
		input.PostalCode,
		input.ReferencePoint,
		input.DeliveryNotes,
		input.IsDefault,
		u.db,
	)
	if err != nil {
		return nil, err
	}

	createdAddress, err := u.addressRepo.CreateAddress(address)
	if err != nil {
		return nil, err
	}

	return createdAddress, nil
}
