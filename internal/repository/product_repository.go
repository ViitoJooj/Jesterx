package repository

import "github.com/ViitoJooj/Jesterx/internal/domain"

type ProductRepository interface {
	CreateProduct(p domain.Product) (*domain.Product, error)
	FindProductByID(id, websiteId string) (*domain.Product, error)
	ListProductsByWebsiteID(websiteId string, onlyActive bool) ([]domain.Product, error)
	UpdateProduct(p domain.Product) (*domain.Product, error)
	DeleteProduct(id, websiteId string) error
}
