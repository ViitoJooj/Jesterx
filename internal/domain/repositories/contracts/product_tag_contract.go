package contracts

import (
	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

type ProductTagContract interface {
	CreateProductTag(tag *domain.ProductsTags) (*domain.ProductsTags, error)
	FindProductTagByUUID(uuid string) (*domain.ProductsTags, error)
	FindProductTagsByLabel(label string) ([]*domain.ProductsTags, error)
	GetProductTagsFromProduct(productUUID string) ([]*domain.ProductsTags, error)
	GetProductTags() ([]*domain.ProductsTags, error)
	UpdateProductTagByUUID(uuid string) error
	DeleteProductTagByUUID(uuid string) error
	DeleteProductTagsByUUIDS(uuid []string) error
}
