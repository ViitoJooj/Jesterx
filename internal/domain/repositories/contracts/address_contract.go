package contracts

import (
	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
)

type AddressContract interface {
	CreateAddress(address *domain.AddressBR) (*domain.AddressBR, error)
	FindAddressByUUID(uuid string) (*domain.AddressBR, error)
	FindDefaultAddressByOwner(ownerUUID string, websiteUUID string) (*domain.AddressBR, error)
	GetAddressesFromOwner(ownerUUID string, websiteUUID string) ([]*domain.AddressBR, error)
	GetAddressesFromWebsite(websiteUUID string) ([]*domain.AddressBR, error)
	GetAddresses() ([]*domain.AddressBR, error)
	UpdateAddressByUUID(uuid string) error
	DeleteAddressByUUID(uuid string) error
	DeleteAddressesByUUIDS(uuid []string) error
}
