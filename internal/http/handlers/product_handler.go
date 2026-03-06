package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: s}
}

type CreateProductRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Price        float64  `json:"price"`
	ComparePrice *float64 `json:"compare_price"`
	Stock        int      `json:"stock"`
	Sku          *string  `json:"sku"`
	Category     *string  `json:"category"`
	Images       []string `json:"images"`
	Active       bool     `json:"active"`
}

type UpdateProductRequest = CreateProductRequest

type ProductData struct {
	ID           string   `json:"id"`
	WebsiteID    string   `json:"website_id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Price        float64  `json:"price"`
	ComparePrice *float64 `json:"compare_price"`
	Stock        int      `json:"stock"`
	Sku          *string  `json:"sku"`
	Category     *string  `json:"category"`
	Images       []string `json:"images"`
	Active       bool     `json:"active"`
	CreatedBy    string   `json:"created_by"`
	UpdatedAt    string   `json:"updated_at"`
	CreatedAt    string   `json:"created_at"`
}

type ProductResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    ProductData `json:"data"`
}

type ProductsResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []ProductData `json:"data"`
}

func productToData(p *domain.Product) ProductData {
	imgs := p.Images
	if imgs == nil {
		imgs = []string{}
	}
	return ProductData{
		ID:           p.Id,
		WebsiteID:    p.WebsiteId,
		Name:         p.Name,
		Description:  p.Description,
		Price:        p.Price,
		ComparePrice: p.ComparePrice,
		Stock:        p.Stock,
		Sku:          p.Sku,
		Category:     p.Category,
		Images:       imgs,
		Active:       p.Active,
		CreatedBy:    p.CreatedBy,
		UpdatedAt:    p.UpdatedAt.Format(time.RFC3339),
		CreatedAt:    p.CreatedAt.Format(time.RFC3339),
	}
}

func productErrStatus(err error) int {
	msg := err.Error()
	if strings.Contains(msg, "acesso negado") || strings.Contains(msg, "não é uma loja") {
		return http.StatusForbidden
	}
	if strings.Contains(msg, "não encontrado") || strings.Contains(msg, "não encontrada") {
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}

// ── Private (owner/admin) ──────────────────────────────────────────────────

// POST /api/v1/sites/{siteID}/products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	defer r.Body.Close()
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p, err := h.productService.CreateProduct(userID, siteID, service.CreateProductInput{
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		ComparePrice: req.ComparePrice,
		Stock:        req.Stock,
		Sku:          req.Sku,
		Category:     req.Category,
		Images:       req.Images,
		Active:       req.Active,
	})
	if err != nil {
		http.Error(w, err.Error(), productErrStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ProductResponse{Success: true, Message: "produto criado", Data: productToData(p)})
}

// GET /api/v1/sites/{siteID}/products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	siteID := strings.TrimSpace(r.PathValue("siteID"))

	products, err := h.productService.ListProducts(userID, siteID)
	if err != nil {
		http.Error(w, err.Error(), productErrStatus(err))
		return
	}

	data := make([]ProductData, 0, len(products))
	for i := range products {
		data = append(data, productToData(&products[i]))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProductsResponse{Success: true, Message: "success", Data: data})
}

// PATCH /api/v1/sites/{siteID}/products/{productID}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	productID := strings.TrimSpace(r.PathValue("productID"))
	defer r.Body.Close()
	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p, err := h.productService.UpdateProduct(userID, siteID, productID, service.UpdateProductInput{
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		ComparePrice: req.ComparePrice,
		Stock:        req.Stock,
		Sku:          req.Sku,
		Category:     req.Category,
		Images:       req.Images,
		Active:       req.Active,
	})
	if err != nil {
		http.Error(w, err.Error(), productErrStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProductResponse{Success: true, Message: "produto atualizado", Data: productToData(p)})
}

// DELETE /api/v1/sites/{siteID}/products/{productID}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	productID := strings.TrimSpace(r.PathValue("productID"))

	if err := h.productService.DeleteProduct(userID, siteID, productID); err != nil {
		http.Error(w, err.Error(), productErrStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Public (store API) ─────────────────────────────────────────────────────

// GET /api/store/{siteID}/products
func (h *ProductHandler) PublicListProducts(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))

	products, err := h.productService.GetPublicProducts(siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := make([]ProductData, 0, len(products))
	for i := range products {
		data = append(data, productToData(&products[i]))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProductsResponse{Success: true, Message: "success", Data: data})
}

// GET /api/store/{siteID}/products/{productID}
func (h *ProductHandler) PublicGetProduct(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	productID := strings.TrimSpace(r.PathValue("productID"))

	p, err := h.productService.GetPublicProduct(siteID, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProductResponse{Success: true, Message: "success", Data: productToData(p)})
}
