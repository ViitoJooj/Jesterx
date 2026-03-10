package service

import (
	"errors"
	"strings"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
)

type ProductService struct {
	productRepo repository.ProductRepository
	websiteRepo repository.WebsiteRepository
	userRepo    repository.UserRepository
}

func NewProductService(
	productRepo repository.ProductRepository,
	websiteRepo repository.WebsiteRepository,
	userRepo repository.UserRepository,
) *ProductService {
	return &ProductService{productRepo, websiteRepo, userRepo}
}

type CreateProductInput struct {
	Name         string
	Description  string
	Price        float64
	ComparePrice *float64
	Stock        int
	Sku          *string
	Category     *string
	Images       []string
	Active       bool
}

type UpdateProductInput struct {
	Name         string
	Description  string
	Price        float64
	ComparePrice *float64
	Stock        int
	Sku          *string
	Category     *string
	Images       []string
	Active       bool
}

// ensures the website exists, is an ECOMMERCE store, and that the user owns it or is an admin
func (s *ProductService) ensureCanManage(userID, websiteID string) (*domain.WebSite, error) {
	website, err := s.websiteRepo.FindWebSiteByID(websiteID)
	if err != nil || website == nil {
		return nil, errors.New("loja não encontrada")
	}
	if website.Banned {
		return nil, errors.New("loja banida")
	}
	if website.Type != "ECOMMERCE" {
		return nil, errors.New("o site não é uma loja ECOMMERCE")
	}

	user, err := s.userRepo.FindUserByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("usuário não encontrado")
	}
	if user.Role != "admin" && website.Creator_id != userID {
		return nil, errors.New("acesso negado: somente o dono da loja ou admin pode gerenciar produtos")
	}
	return website, nil
}

func (s *ProductService) CreateProduct(userID, websiteID string, input CreateProductInput) (*domain.Product, error) {
	if _, err := s.ensureCanManage(userID, websiteID); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(input.Name)
	if len(name) < 2 || len(name) > 200 {
		return nil, errors.New("nome do produto inválido (2-200 caracteres)")
	}
	if input.Price < 0 {
		return nil, errors.New("preço não pode ser negativo")
	}
	if input.Stock < 0 {
		return nil, errors.New("estoque não pode ser negativo")
	}
	if input.Images == nil {
		input.Images = []string{}
	}

	p := domain.NewProduct(
		websiteID, name, input.Description,
		input.Price, input.ComparePrice,
		input.Stock, input.Sku, input.Category,
		input.Images, input.Active, userID,
	)
	return s.productRepo.CreateProduct(*p)
}

func (s *ProductService) ListProducts(userID, websiteID string) ([]domain.Product, error) {
	if _, err := s.ensureCanManage(userID, websiteID); err != nil {
		return nil, err
	}
	return s.productRepo.ListProductsByWebsiteID(websiteID, false)
}

func (s *ProductService) GetPublicProducts(websiteID string) ([]domain.Product, error) {
	website, err := s.websiteRepo.FindWebSiteByID(websiteID)
	if err != nil || website == nil {
		return nil, errors.New("loja não encontrada")
	}
	if website.Banned {
		return nil, errors.New("loja banida")
	}
	if website.Type != "ECOMMERCE" {
		return nil, errors.New("o site não é uma loja")
	}
	return s.productRepo.ListProductsByWebsiteID(websiteID, true)
}

func (s *ProductService) GetPublicProduct(websiteID, productID string) (*domain.Product, error) {
	website, err := s.websiteRepo.FindWebSiteByID(websiteID)
	if err != nil || website == nil {
		return nil, errors.New("loja não encontrada")
	}
	if website.Banned || website.Type != "ECOMMERCE" {
		return nil, errors.New("loja inválida")
	}
	p, err := s.productRepo.FindProductByID(productID, websiteID)
	if err != nil || p == nil {
		return nil, errors.New("produto não encontrado")
	}
	if !p.Active {
		return nil, errors.New("produto não disponível")
	}
	return p, nil
}

func (s *ProductService) UpdateProduct(userID, websiteID, productID string, input UpdateProductInput) (*domain.Product, error) {
	if _, err := s.ensureCanManage(userID, websiteID); err != nil {
		return nil, err
	}

	existing, err := s.productRepo.FindProductByID(productID, websiteID)
	if err != nil || existing == nil {
		return nil, errors.New("produto não encontrado")
	}

	name := strings.TrimSpace(input.Name)
	if len(name) < 2 || len(name) > 200 {
		return nil, errors.New("nome do produto inválido (2-200 caracteres)")
	}
	if input.Price < 0 {
		return nil, errors.New("preço não pode ser negativo")
	}
	if input.Stock < 0 {
		return nil, errors.New("estoque não pode ser negativo")
	}
	if input.Images == nil {
		input.Images = existing.Images
	}

	existing.Name = name
	existing.Description = input.Description
	existing.Price = input.Price
	existing.ComparePrice = input.ComparePrice
	existing.Stock = input.Stock
	existing.Sku = input.Sku
	existing.Category = input.Category
	existing.Images = input.Images
	existing.Active = input.Active

	return s.productRepo.UpdateProduct(*existing)
}

func (s *ProductService) DeleteProduct(userID, websiteID, productID string) error {
	if _, err := s.ensureCanManage(userID, websiteID); err != nil {
		return err
	}
	return s.productRepo.DeleteProduct(productID, websiteID)
}
