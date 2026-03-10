package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
)

type CreateOrderInput struct {
	BuyerName  string
	BuyerEmail string
	BuyerPhone string
	Notes      string
	Items      []OrderItemInput
}

type OrderItemInput struct {
	ProductID string
	Qty       int
}

type OrderService struct {
	orderRepo   repository.OrderRepository
	websiteRepo repository.WebsiteRepository
	productRepo repository.ProductRepository
}

func NewOrderService(or repository.OrderRepository, wr repository.WebsiteRepository, pr repository.ProductRepository) *OrderService {
	return &OrderService{orderRepo: or, websiteRepo: wr, productRepo: pr}
}

func platformCommissionPct() float64 {
	v := os.Getenv("PLATFORM_COMMISSION_PCT")
	if v == "" {
		return 5.0
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || f < 0 || f > 100 {
		return 5.0
	}
	return f
}

// ensureSiteOrderable validates that the site exists, is not banned, and is ECOMMERCE type.
func (s *OrderService) ensureSiteOrderable(websiteID string) (*domain.WebSite, error) {
	site, err := s.websiteRepo.FindWebSiteByID(websiteID)
	if err != nil {
		return nil, errors.New("erro interno")
	}
	if site == nil {
		return nil, errors.New("loja não encontrada")
	}
	if site.Banned {
		return nil, errors.New("loja indisponível")
	}
	if site.Type != "ECOMMERCE" {
		return nil, errors.New("loja não aceita pedidos")
	}
	return site, nil
}

// CreateOrder creates a new order. Product prices are always read from the
// database — the client only provides product IDs and quantities. This prevents
// price manipulation attacks.
func (s *OrderService) CreateOrder(websiteID string, input CreateOrderInput) (*domain.Order, error) {
	if len(input.Items) == 0 {
		return nil, fmt.Errorf("pedido deve ter pelo menos 1 item")
	}
	if input.BuyerEmail == "" {
		return nil, fmt.Errorf("email do comprador é obrigatório")
	}

	if _, err := s.ensureSiteOrderable(websiteID); err != nil {
		return nil, err
	}

	var subtotal float64
	items := make([]domain.OrderItem, 0, len(input.Items))
	for _, it := range input.Items {
		if it.ProductID == "" {
			return nil, fmt.Errorf("product_id é obrigatório")
		}
		if it.Qty <= 0 {
			it.Qty = 1
		}

		// Always fetch authoritative price from DB to prevent price manipulation.
		product, err := s.productRepo.FindProductByID(it.ProductID, websiteID)
		if err != nil {
			return nil, errors.New("erro ao buscar produto")
		}
		if product == nil || !product.Active {
			return nil, fmt.Errorf("produto '%s' não encontrado ou inativo", it.ProductID)
		}
		if product.Stock < it.Qty {
			return nil, fmt.Errorf("estoque insuficiente para '%s'", product.Name)
		}

		lineTotal := product.Price * float64(it.Qty)
		subtotal += lineTotal
		items = append(items, domain.OrderItem{
			ProductID:   it.ProductID,
			ProductName: product.Name,
			UnitPrice:   product.Price,
			Qty:         it.Qty,
			Total:       lineTotal,
		})
	}

	pct := platformCommissionPct()
	fee := subtotal * (pct / 100.0)

	order := &domain.Order{
		WebsiteID:   websiteID,
		BuyerName:   input.BuyerName,
		BuyerEmail:  input.BuyerEmail,
		BuyerPhone:  input.BuyerPhone,
		Status:      domain.OrderPending,
		Subtotal:    subtotal,
		PlatformFee: fee,
		Total:       subtotal,
		Notes:       input.Notes,
		Items:       items,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("criar pedido: %w", err)
	}

	return order, nil
}

// ListSiteOrders returns orders for a site after verifying ownership.
func (s *OrderService) ListSiteOrders(userID, websiteID string) ([]domain.Order, error) {
	site, err := s.websiteRepo.FindWebSiteByID(websiteID)
	if err != nil {
		return nil, errors.New("erro interno")
	}
	if site == nil {
		return nil, errors.New("site não encontrado")
	}
	if site.Creator_id != userID {
		return nil, errors.New("acesso negado")
	}
	return s.orderRepo.ListBySite(websiteID)
}

func (s *OrderService) GetSiteOrdersSince(from, to time.Time) ([]domain.Order, error) {
	return s.orderRepo.ListSince(from, to)
}
