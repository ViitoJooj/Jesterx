package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: s}
}

type CreateOrderRequest struct {
	BuyerName  string         `json:"buyer_name"`
	BuyerEmail string         `json:"buyer_email"`
	BuyerPhone string         `json:"buyer_phone"`
	Notes      string         `json:"notes"`
	Items      []OrderItemReq `json:"items"`
}

type OrderItemReq struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	UnitPrice   float64 `json:"unit_price"`
	Qty         int     `json:"qty"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	defer r.Body.Close()
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	items := make([]service.OrderItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, service.OrderItemInput{
			ProductID:   it.ProductID,
			ProductName: it.ProductName,
			UnitPrice:   it.UnitPrice,
			Qty:         it.Qty,
		})
	}

	order, err := h.orderService.CreateOrder(siteID, service.CreateOrderInput{
		BuyerName:  req.BuyerName,
		BuyerEmail: req.BuyerEmail,
		BuyerPhone: req.BuyerPhone,
		Notes:      req.Notes,
		Items:      items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "pedido criado", "data": order})
}

func (h *OrderHandler) ListSiteOrders(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	siteID := strings.TrimSpace(r.PathValue("siteID"))

	orders, err := h.orderService.ListSiteOrders(siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "success", "data": orders})
}
