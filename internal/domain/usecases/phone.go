package usecases

import (

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/google/uuid"
)

type CreatePhoneUseCase struct {
	phoneRepo contracts.PhoneContract
}

func NewCreatePhoneUseCase(phoneRepo contracts.PhoneContract) *CreatePhoneUseCase {
	return &CreatePhoneUseCase{
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

func (u *CreatePhoneUseCase) GetByUUID(uuidStr string) (*domain.Phone, error) {
	return u.phoneRepo.FindPhoneByUUID(uuidStr)
}

func (u *CreatePhoneUseCase) GetAll(websiteUUIDStr string) ([]*domain.Phone, error) {
	return u.phoneRepo.GetPhonesFromWebsite(websiteUUIDStr)
}
