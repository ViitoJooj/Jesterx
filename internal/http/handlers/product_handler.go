package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/service"
	"github.com/ViitoJooj/Jesterx/pkg/validate"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: s}
}

type CreateProductRequest struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	ShortDescription *string           `json:"short_description"`
	Price            float64           `json:"price"`
	ComparePrice     *float64          `json:"compare_price"`
	Stock            int               `json:"stock"`
	Sku              *string           `json:"sku"`
	Category         *string           `json:"category"`
	Slug             *string           `json:"slug"`
	Brand            *string           `json:"brand"`
	Model            *string           `json:"model"`
	Barcode          *string           `json:"barcode"`
	Condition        *string           `json:"condition"`
	WeightGrams      *int              `json:"weight_grams"`
	WidthCm          *float64          `json:"width_cm"`
	HeightCm         *float64          `json:"height_cm"`
	LengthCm         *float64          `json:"length_cm"`
	Material         *string           `json:"material"`
	Color            *string           `json:"color"`
	Size             *string           `json:"size"`
	WarrantyMonths   *int              `json:"warranty_months"`
	OriginCountry    *string           `json:"origin_country"`
	Tags             []string          `json:"tags"`
	Attributes       map[string]string `json:"attributes"`
	RequiresShipping *bool             `json:"requires_shipping"`
	Images           []string          `json:"images"`
	Active           bool              `json:"active"`
}

type UpdateProductRequest = CreateProductRequest

type ProductData struct {
	ID               string            `json:"id"`
	WebsiteID        string            `json:"website_id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	ShortDescription *string           `json:"short_description,omitempty"`
	Price            float64           `json:"price"`
	ComparePrice     *float64          `json:"compare_price,omitempty"`
	Stock            int               `json:"stock"`
	Sku              *string           `json:"sku,omitempty"`
	Category         *string           `json:"category,omitempty"`
	Slug             *string           `json:"slug,omitempty"`
	Brand            *string           `json:"brand,omitempty"`
	Model            *string           `json:"model,omitempty"`
	Barcode          *string           `json:"barcode,omitempty"`
	Condition        *string           `json:"condition,omitempty"`
	WeightGrams      *int              `json:"weight_grams,omitempty"`
	WidthCm          *float64          `json:"width_cm,omitempty"`
	HeightCm         *float64          `json:"height_cm,omitempty"`
	LengthCm         *float64          `json:"length_cm,omitempty"`
	Material         *string           `json:"material,omitempty"`
	Color            *string           `json:"color,omitempty"`
	Size             *string           `json:"size,omitempty"`
	WarrantyMonths   *int              `json:"warranty_months,omitempty"`
	OriginCountry    *string           `json:"origin_country,omitempty"`
	Tags             []string          `json:"tags"`
	Attributes       map[string]string `json:"attributes"`
	RequiresShipping bool              `json:"requires_shipping"`
	Images           []string          `json:"images"`
	Active           bool              `json:"active"`
	CreatedBy        string            `json:"created_by"`
	UpdatedAt        string            `json:"updated_at"`
	CreatedAt        string            `json:"created_at"`
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
		ID:               p.Id,
		WebsiteID:        p.WebsiteId,
		Name:             p.Name,
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		Price:            p.Price,
		ComparePrice:     p.ComparePrice,
		Stock:            p.Stock,
		Sku:              p.Sku,
		Category:         p.Category,
		Slug:             p.Slug,
		Brand:            p.Brand,
		Model:            p.Model,
		Barcode:          p.Barcode,
		Condition:        p.Condition,
		WeightGrams:      p.WeightGrams,
		WidthCm:          p.WidthCm,
		HeightCm:         p.HeightCm,
		LengthCm:         p.LengthCm,
		Material:         p.Material,
		Color:            p.Color,
		Size:             p.Size,
		WarrantyMonths:   p.WarrantyMonths,
		OriginCountry:    p.OriginCountry,
		Tags:             p.Tags,
		Attributes:       p.Attributes,
		RequiresShipping: p.RequiresShipping,
		Images:           imgs,
		Active:           p.Active,
		CreatedBy:        p.CreatedBy,
		UpdatedAt:        p.UpdatedAt.Format(time.RFC3339),
		CreatedAt:        p.CreatedAt.Format(time.RFC3339),
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

	if err := validate.New().
		Required("name", req.Name).
		MaxLen("name", req.Name, 200).
		MinFloat("price", req.Price, 0.01).
		MinInt("stock", req.Stock, 0).
		Err(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	p, err := h.productService.CreateProduct(userID, siteID, service.CreateProductInput{
		Name:             req.Name,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Price:            req.Price,
		ComparePrice:     req.ComparePrice,
		Stock:            req.Stock,
		Sku:              req.Sku,
		Category:         req.Category,
		Slug:             req.Slug,
		Brand:            req.Brand,
		Model:            req.Model,
		Barcode:          req.Barcode,
		Condition:        req.Condition,
		WeightGrams:      req.WeightGrams,
		WidthCm:          req.WidthCm,
		HeightCm:         req.HeightCm,
		LengthCm:         req.LengthCm,
		Material:         req.Material,
		Color:            req.Color,
		Size:             req.Size,
		WarrantyMonths:   req.WarrantyMonths,
		OriginCountry:    req.OriginCountry,
		Tags:             req.Tags,
		Attributes:       req.Attributes,
		RequiresShipping: req.RequiresShipping,
		Images:           req.Images,
		Active:           req.Active,
	})
	if err != nil {
		http.Error(w, err.Error(), productErrStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ProductResponse{Success: true, Message: "produto criado", Data: productToData(p)})
}

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

	if err := validate.New().
		Required("name", req.Name).
		MaxLen("name", req.Name, 200).
		MinFloat("price", req.Price, 0.01).
		MinInt("stock", req.Stock, 0).
		Err(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	p, err := h.productService.UpdateProduct(userID, siteID, productID, service.UpdateProductInput{
		Name:             req.Name,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Price:            req.Price,
		ComparePrice:     req.ComparePrice,
		Stock:            req.Stock,
		Sku:              req.Sku,
		Category:         req.Category,
		Slug:             req.Slug,
		Brand:            req.Brand,
		Model:            req.Model,
		Barcode:          req.Barcode,
		Condition:        req.Condition,
		WeightGrams:      req.WeightGrams,
		WidthCm:          req.WidthCm,
		HeightCm:         req.HeightCm,
		LengthCm:         req.LengthCm,
		Material:         req.Material,
		Color:            req.Color,
		Size:             req.Size,
		WarrantyMonths:   req.WarrantyMonths,
		OriginCountry:    req.OriginCountry,
		Tags:             req.Tags,
		Attributes:       req.Attributes,
		RequiresShipping: req.RequiresShipping,
		Images:           req.Images,
		Active:           req.Active,
	})
	if err != nil {
		http.Error(w, err.Error(), productErrStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProductResponse{Success: true, Message: "produto atualizado", Data: productToData(p)})
}

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
