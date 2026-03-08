package service

import (
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
	ProductID   string
	ProductName string
	UnitPrice   float64
	Qty         int
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

func (s *OrderService) CreateOrder(websiteID string, input CreateOrderInput) (*domain.Order, error) {
	if len(input.Items) == 0 {
		return nil, fmt.Errorf("pedido deve ter pelo menos 1 item")
	}
	if input.BuyerEmail == "" {
		return nil, fmt.Errorf("email do comprador é obrigatório")
	}

	var subtotal float64
	items := make([]domain.OrderItem, 0, len(input.Items))
	for _, it := range input.Items {
		if it.Qty <= 0 {
			it.Qty = 1
		}
		lineTotal := it.UnitPrice * float64(it.Qty)
		subtotal += lineTotal
		items = append(items, domain.OrderItem{
			ProductID:   it.ProductID,
			ProductName: it.ProductName,
			UnitPrice:   it.UnitPrice,
			Qty:         it.Qty,
			Total:       lineTotal,
		})
	}

	pct := platformCommissionPct()
	fee := subtotal * (pct / 100.0)
	total := subtotal

	order := &domain.Order{
		WebsiteID:   websiteID,
		BuyerName:   input.BuyerName,
		BuyerEmail:  input.BuyerEmail,
		BuyerPhone:  input.BuyerPhone,
		Status:      domain.OrderPending,
		Subtotal:    subtotal,
		PlatformFee: fee,
		Total:       total,
		Notes:       input.Notes,
		Items:       items,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("criar pedido: %w", err)
	}

	return order, nil
}

func (s *OrderService) ListSiteOrders(websiteID string) ([]domain.Order, error) {
	return s.orderRepo.ListBySite(websiteID)
}

func (s *OrderService) GetSiteOrdersSince(from, to time.Time) ([]domain.Order, error) {
	return s.orderRepo.ListSince(from, to)
}
