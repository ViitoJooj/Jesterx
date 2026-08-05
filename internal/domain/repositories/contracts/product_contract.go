package contracts

import (
	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

type ProductContract interface {
	CreateProduct(product *domain.Products) (*domain.Products, error)
	FindProductByUUID(uuid string) (*domain.Products, error)
	FindProductByName(name string) (*domain.Products, error)
	GetProducts() ([]*domain.Products, error)
	GetActiveProducts() ([]*domain.Products, error)
	UpdateProductByUUID(uuid string) error
	DeleteProductByUUID(uuid string) error
	DeleteProductsByUUIDS(uuid []string) error
}
