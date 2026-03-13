package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
	"github.com/ViitoJooj/Jesterx/internal/security"
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
	userRepo    repository.UserRepository
}

func NewOrderService(or repository.OrderRepository, wr repository.WebsiteRepository, pr repository.ProductRepository, ur repository.UserRepository) *OrderService {
	return &OrderService{orderRepo: or, websiteRepo: wr, productRepo: pr, userRepo: ur}
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
func (s *OrderService) CreateOrder(userID, websiteID string, input CreateOrderInput) (*domain.Order, error) {
	if len(input.Items) == 0 {
		return nil, fmt.Errorf("pedido deve ter pelo menos 1 item")
	}

	site, err := s.ensureSiteOrderable(websiteID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindUserByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("usuário não encontrado")
	}
	if user.WebsiteId != websiteID {
		return nil, errors.New("usuário inválido para esta loja")
	}
	if missingAddress(user) {
		return nil, errors.New("endereço obrigatório para finalizar a compra")
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

	buyerName := strings.TrimSpace(input.BuyerName)
	if buyerName == "" {
		buyerName = strings.TrimSpace(user.First_name + " " + user.Last_name)
	}
	if buyerName == "" {
		buyerName = user.Email
	}

	buyerPhone := ""
	if user.Phone != nil {
		buyerPhone = strings.TrimSpace(*user.Phone)
	}
	if strings.TrimSpace(input.BuyerPhone) != "" {
		buyerPhone = strings.TrimSpace(input.BuyerPhone)
	}

	total := subtotal
	order := &domain.Order{
		WebsiteID:                 websiteID,
		BuyerUserID:               &userID,
		BuyerName:                 buyerName,
		BuyerEmail:                user.Email,
		BuyerPhone:                buyerPhone,
		BuyerDocument:             derefString(user.CpfCnpj),
		ShippingName:              buyerName,
		ShippingPhone:             buyerPhone,
		ShippingZipCode:           derefString(user.ZipCode),
		ShippingAddressStreet:     derefString(user.AddressStreet),
		ShippingAddressNumber:     derefString(user.AddressNumber),
		ShippingAddressComplement: derefString(user.AddressComplement),
		ShippingAddressDistrict:   derefString(user.AddressDistrict),
		ShippingAddressCity:       derefString(user.AddressCity),
		ShippingAddressState:      derefString(user.AddressState),
		ShippingAddressCountry:    derefString(user.AddressCountry),
		ShippingCost:              0,
		DiscountTotal:             0,
		TaxTotal:                  0,
		Currency:                  "BRL",
		Status:                    domain.OrderPending,
		Subtotal:                  subtotal,
		PlatformFee:               fee,
		Total:                     total,
		Notes:                     input.Notes,
		Items:                     items,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("criar pedido: %w", err)
	}

	s.notifyOwner(site, order)

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

func (s *OrderService) notifyOwner(site *domain.WebSite, order *domain.Order) {
	if site == nil {
		return
	}
	owner, err := s.userRepo.FindUserByID(site.Creator_id)
	if err != nil || owner == nil || owner.Email == "" {
		return
	}
	subject := fmt.Sprintf("🛍 Novo pedido em %s", site.Name)
	body := buildOrderNotificationEmail(owner.First_name, site.Name, order)
	if err := security.SendOrderNotificationEmail(owner.Email, subject, body); err != nil {
		log.Printf("[order_notify] erro ao enviar para %s: %v", owner.Email, err)
	}
}

func buildOrderNotificationEmail(ownerName, siteName string, order *domain.Order) string {
	var items strings.Builder
	for _, it := range order.Items {
		items.WriteString(fmt.Sprintf("<li>%dx %s — R$ %.2f</li>", it.Qty, it.ProductName, it.Total))
	}
	address := strings.TrimSpace(fmt.Sprintf("%s, %s", order.ShippingAddressStreet, order.ShippingAddressNumber))
	if order.ShippingAddressComplement != "" {
		address = address + " - " + order.ShippingAddressComplement
	}
	address = strings.TrimSpace(address)
	location := strings.TrimSpace(fmt.Sprintf("%s - %s", order.ShippingAddressCity, order.ShippingAddressState))

	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="utf-8"/></head>
<body style="margin:0;font-family:Inter,system-ui,sans-serif;background:#f5f7fa">
<div style="max-width:600px;margin:32px auto;background:#fff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,.08)">
  <div style="background:linear-gradient(135deg,#1a2740,#2d4070);padding:28px 32px;color:#fff">
    <h1 style="margin:0;font-size:20px">🛍 Novo pedido</h1>
    <p style="margin:8px 0 0;opacity:.8;font-size:13px">%s</p>
  </div>
  <div style="padding:28px 32px">
    <p style="font-size:15px;color:#1a2740">Olá, <strong>%s</strong>! Você recebeu um novo pedido.</p>
    <div style="background:#f5f7fa;border-radius:10px;padding:14px 16px;margin:16px 0">
      <div style="font-size:13px;color:#6b7280">Comprador</div>
      <div style="font-size:15px;font-weight:600;color:#1a2740">%s</div>
      <div style="font-size:13px;color:#6b7280">%s</div>
    </div>
    <div style="font-size:13px;color:#6b7280;margin-top:10px">Itens</div>
    <ul style="margin:8px 0 16px;padding-left:18px;font-size:14px;color:#1a2740">%s</ul>
    <div style="display:flex;justify-content:space-between;align-items:center;background:#f5f7fa;border-radius:10px;padding:12px 16px">
      <span style="font-size:13px;color:#6b7280">Total</span>
      <span style="font-size:18px;font-weight:700;color:#ff5d1f">R$ %.2f</span>
    </div>
    <div style="margin-top:16px;font-size:13px;color:#6b7280">
      <div><strong>Entrega:</strong> %s</div>
      <div>%s</div>
      <div>%s</div>
    </div>
  </div>
  <div style="padding:18px 32px;border-top:1px solid #f0f0f0;text-align:center;font-size:12px;color:#9aa5bc">
    Jesterx · %s
  </div>
</div>
</body></html>`,
		siteName,
		ownerName,
		order.BuyerName,
		order.BuyerEmail,
		items.String(),
		order.Total,
		order.ShippingZipCode,
		address,
		location,
		siteName,
	)
}

func missingAddress(user *domain.User) bool {
	return user.ZipCode == nil || strings.TrimSpace(*user.ZipCode) == "" ||
		user.AddressStreet == nil || strings.TrimSpace(*user.AddressStreet) == "" ||
		user.AddressNumber == nil || strings.TrimSpace(*user.AddressNumber) == "" ||
		user.AddressCity == nil || strings.TrimSpace(*user.AddressCity) == "" ||
		user.AddressState == nil || strings.TrimSpace(*user.AddressState) == "" ||
		user.AddressCountry == nil || strings.TrimSpace(*user.AddressCountry) == ""
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}
