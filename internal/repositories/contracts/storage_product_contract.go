package contracts

import (
	"github.com/ViitoJooj/Jesterx/internal/domain"
)

type StorageProductContract interface {
	CreateStorageProduct(sp *domain.StorageProducts) (*domain.StorageProducts, error)
	FindStorageProductByUUID(uuid string) (*domain.StorageProducts, error)
	FindStorageProductByProductUUID(productUUID string) (*domain.StorageProducts, error)
	GetStorageProducts() ([]*domain.StorageProducts, error)
	UpdateStorageProductByUUID(uuid string) error
	DeleteStorageProductByUUID(uuid string) error
	DeleteStorageProductsByUUIDS(uuid []string) error
}
