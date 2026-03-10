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

// OrderItemReq only accepts product_id and qty — prices are always
// resolved server-side from the database to prevent price manipulation.
type OrderItemReq struct {
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	defer r.Body.Close()
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"message":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	items := make([]service.OrderItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, service.OrderItemInput{
			ProductID: it.ProductID,
			Qty:       it.Qty,
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "message": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "pedido criado", "data": order})
}

func (h *OrderHandler) ListSiteOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, `{"success":false,"message":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	siteID := strings.TrimSpace(r.PathValue("siteID"))

	orders, err := h.orderService.ListSiteOrders(userID, siteID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "acesso negado" {
			status = http.StatusForbidden
		} else if err.Error() == "site não encontrado" {
			status = http.StatusNotFound
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "message": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "success", "data": orders})
}
