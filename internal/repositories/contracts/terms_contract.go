package contracts

import (
	"github.com/ViitoJooj/Jesterx/internal/domain"
)

type TermsContract interface {
	CreateTerms(terms *domain.Terms) (*domain.Terms, error)
	FindTermsByUUID(uuid string) (*domain.Terms, error)
	FindTermsByName(name string) (*domain.Terms, error)
	GetTerms() ([]*domain.Terms, error)
	UpdateTermsByUUID(uuid string) error
	DeleteTermsByUUID(uuid string) error
	DeleteTermsByUUIDS(uuid []string) error
}
