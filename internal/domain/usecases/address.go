package usecases

import (
	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/google/uuid"
)

type CreateAddressUseCase struct {
	addressRepo contracts.AddressContract
}

func NewCreateAddressUseCase(addressRepo contracts.AddressContract) *CreateAddressUseCase {
	return &CreateAddressUseCase{
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

func (u *CreateAddressUseCase) GetByUUID(uuidStr string) (*domain.AddressBR, error) {
	return u.addressRepo.FindAddressByUUID(uuidStr)
}

func (u *CreateAddressUseCase) GetAll(websiteUUIDStr string) ([]*domain.AddressBR, error) {
	return u.addressRepo.GetAddressesFromWebsite(websiteUUIDStr)
}
