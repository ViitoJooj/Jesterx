package contracts

import (
	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
)

type CupomContract interface {
	CreateCupom(cupom *domain.Cupons) (*domain.Cupons, error)
	FindCupomByUUID(uuid string) (*domain.Cupons, error)
	FindCupomByLabel(label string) (*domain.Cupons, error)
	GetCuponsFromTag(tagUUID string) ([]*domain.Cupons, error)
	GetCupons() ([]*domain.Cupons, error)
	UpdateCupomByUUID(uuid string) error
	DeleteCupomByUUID(uuid string) error
	DeleteCupomsByUUIDS(uuid []string) error
}
