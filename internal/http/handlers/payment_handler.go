package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

type PlanResponse struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	DescriptionM string  `json:"description_md"`
	Price        float64 `json:"price"`
	BillingCycle string  `json:"billing_cycle"`
}

type ListPlansResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []PlanResponse `json:"data"`
}

type CreateCheckoutRequest struct {
	PlanID   int64 `json:"plan_id"`
	Quantity int   `json:"quantity"`
}

type CheckoutData struct {
	SessionID   string `json:"session_id"`
	CheckoutURL string `json:"checkout_url"`
}

type CheckoutResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    CheckoutData `json:"data"`
}

func (h *PaymentHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.paymentService.ListActivePlans()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	respPlans := make([]PlanResponse, 0, len(plans))
	for _, p := range plans {
		respPlans = append(respPlans, PlanResponse{
			ID:           p.ID,
			Name:         p.Name,
			Description:  p.Description,
			DescriptionM: p.DescriptionM,
			Price:        p.Price,
			BillingCycle: p.BillingCycle,
		})
	}

	resp := ListPlansResponse{
		Success: true,
		Message: "success",
		Data:    respPlans,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *PaymentHandler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()
	var req CreateCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.paymentService.CreateCheckoutSession(userID, req.PlanID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := CheckoutResponse{
		Success: true,
		Message: "checkout created",
		Data: CheckoutData{
			SessionID:   result.SessionID,
			CheckoutURL: result.CheckoutURL,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *PaymentHandler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("Stripe-Signature")
	if err := h.paymentService.ProcessStripeWebhook(body, signature); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
