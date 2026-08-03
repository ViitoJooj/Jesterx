package contracts

import (
	"github.com/ViitoJooj/Jesterx/internal/domain"
)

type PhoneContract interface {
	CreatePhone(phone *domain.Phone) (*domain.Phone, error)
	FindPhoneByUUID(uuid string) (*domain.Phone, error)
	FindDefaultPhoneByOwner(ownerUUID string, websiteUUID string) (*domain.Phone, error)
	GetPhonesFromOwner(ownerUUID string, websiteUUID string) ([]*domain.Phone, error)
	GetPhonesFromWebsite(websiteUUID string) ([]*domain.Phone, error)
	GetPhones() ([]*domain.Phone, error)
	UpdatePhoneByUUID(uuid string) error
	DeletePhoneByUUID(uuid string) error
	DeletePhonesByUUIDS(uuid []string) error
}
