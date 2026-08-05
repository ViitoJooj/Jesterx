package contracts

import (
	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

type ProductShippedContract interface {
	CreateProductShipped(productShipped *domain.ProductShipped) (*domain.ProductShipped, error)
	FindProductShippedByUUID(uuid string) (*domain.ProductShipped, error)
	FindProductShippedByProductUUID(productUUID string) ([]*domain.ProductShipped, error)
	FindProductShippedByStatus(status string) ([]*domain.ProductShipped, error)
	GetProductsShipped() ([]*domain.ProductShipped, error)
	UpdateProductShippedByUUID(uuid string) error
	DeleteProductShippedByUUID(uuid string) error
	DeleteProductsShippedByUUIDS(uuid []string) error
}
