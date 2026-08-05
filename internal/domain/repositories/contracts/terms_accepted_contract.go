package contracts

import (
	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

type TermsAcceptedContract interface {
	CreateTermsAccepted(termsAccepted *domain.TermsAcceptedBy) (*domain.TermsAcceptedBy, error)
	FindTermsAcceptedByUUID(uuid string) (*domain.TermsAcceptedBy, error)
	FindTermsAcceptedByOwner(ownerUUID string, websiteUUID string) (*domain.TermsAcceptedBy, error)
	GetTermsAcceptedFromWebsite(websiteUUID string) ([]*domain.TermsAcceptedBy, error)
	GetTermsAccepted() ([]*domain.TermsAcceptedBy, error)
	UpdateTermsAcceptedByUUID(uuid string) error
	DeleteTermsAcceptedByUUID(uuid string) error
	DeleteTermsAcceptedByUUIDS(uuid []string) error
}
