package contracts

import (
	"github.com/ViitoJooj/Jesterx/internal/domain"
)

type PreparingShippingProductContract interface {
	CreatePreparingShippingProduct(psp *domain.PreparingShippingProducts) (*domain.PreparingShippingProducts, error)
	FindPreparingShippingProductByUUID(uuid string) (*domain.PreparingShippingProducts, error)
	FindPreparingShippingProductByProductUUID(productUUID string) (*domain.PreparingShippingProducts, error)
	GetPreparingShippingProducts() ([]*domain.PreparingShippingProducts, error)
	UpdatePreparingShippingProductByUUID(uuid string) error
	DeletePreparingShippingProductByUUID(uuid string) error
	DeletePreparingShippingProductsByUUIDS(uuid []string) error
}
