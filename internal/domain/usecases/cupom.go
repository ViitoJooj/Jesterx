package usecases

import (

	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
)

type CreateCupomUseCase struct {
	repository contracts.CupomContract
}

func NewCreateCupomUseCase(repository contracts.CupomContract) *CreateCupomUseCase {
	return &CreateCupomUseCase{repository: repository}
}

func (u *CreateCupomUseCase) Create(tagUUID string, label string, description string, value string, valueType string) (*domain.Cupons, error) {
	cupom, err := domain.NewCupom(tagUUID, label, description, value, valueType)
	if err != nil {
		return nil, err
	}

	createdCupom, err := u.repository.CreateCupom(cupom)
	if err != nil {
		return nil, err
	}

	return createdCupom, nil
}
