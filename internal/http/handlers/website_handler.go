package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

type WebSiteRequest struct {
	Type              string `json:"type"`
	Image             []byte `json:"image,omitempty"`
	Name              string `json:"name"`
	Short_description string `json:"short_description,omitempty"`
	Description       string `json:"description,omitempty"`
}

type WebsiteRouteRequest struct {
	Path         string `json:"path"`
	Title        string `json:"title"`
	RequiresAuth bool   `json:"requires_auth"`
}

type ReplaceRoutesRequest struct {
	Routes []WebsiteRouteRequest `json:"routes"`
}

type CreateVersionRequest struct {
	SourceType string `json:"source_type"`
	Source     string `json:"source"`
}

type WebSiteData struct {
	Id                string    `json:"id"`
	Type              string    `json:"type"`
	Image             []byte    `json:"image"`
	Name              string    `json:"name"`
	Short_description string    `json:"short_description"`
	Description       string    `json:"description"`
	Creator_id        string    `json:"creator_id"`
	Banned            bool      `json:"banned"`
	Updated_at        time.Time `json:"updated_at"`
	Created_at        time.Time `json:"created_at"`
}

type WebSiteResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    WebSiteData `json:"data"`
}

type WebSitesResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []WebSiteData `json:"data"`
}

type RouteResponse struct {
	ID           string `json:"id"`
	Path         string `json:"path"`
	Title        string `json:"title"`
	RequiresAuth bool   `json:"requires_auth"`
	Position     int    `json:"position"`
}

type RoutesResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    []RouteResponse `json:"data"`
}

type VersionData struct {
	ID           string  `json:"id"`
	WebsiteID    string  `json:"website_id"`
	Version      int     `json:"version"`
	SourceType   string  `json:"source_type"`
	Source       string  `json:"source,omitempty"`
	ScanStatus   string  `json:"scan_status"`
	ScanScore    int     `json:"scan_score"`
	ScanFindings string  `json:"scan_findings"`
	Published    bool    `json:"published"`
	PublishedAt  *string `json:"published_at,omitempty"`
}

type VersionResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    VersionData `json:"data"`
}

type VersionsResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []VersionData `json:"data"`
}

type SiteAPIItem struct {
	ID          string `json:"id"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type SiteAPIsResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []SiteAPIItem `json:"data"`
}

type ScanReportResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Version      int    `json:"version"`
		ScanStatus   string `json:"scan_status"`
		ScanScore    int    `json:"scan_score"`
		ScanFindings string `json:"scan_findings"`
		SourceType   string `json:"source_type"`
	} `json:"data"`
}

type WebSiteHandler struct {
	webSiteService *service.WebSiteService
}

func NewWebSiteHandler(webSiteService *service.WebSiteService) *WebSiteHandler {
	return &WebSiteHandler{
		webSiteService: webSiteService,
	}
}

func asVersionData(versionID string, websiteID string, version int, sourceType string, source string, scanStatus string, scanScore int, scanFindings string, published bool, publishedAt *time.Time) VersionData {
	var publishedAtStr *string
	if publishedAt != nil {
		v := publishedAt.Format(time.RFC3339)
		publishedAtStr = &v
	}
	return VersionData{
		ID:           versionID,
		WebsiteID:    websiteID,
		Version:      version,
		SourceType:   sourceType,
		Source:       source,
		ScanStatus:   scanStatus,
		ScanScore:    scanScore,
		ScanFindings: scanFindings,
		Published:    published,
		PublishedAt:  publishedAtStr,
	}
}

func (h *WebSiteHandler) CreateWebSite(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req WebSiteRequest

	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	website, err := h.webSiteService.CreateWebSite(req.Type, req.Image, req.Name, req.Short_description, req.Description, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := WebSiteResponse{
		Success: true,
		Message: fmt.Sprintf("%s created.", strings.ToUpper(req.Type)),
		Data: WebSiteData{
			Id:                website.Id,
			Type:              website.Type,
			Image:             website.Image,
			Name:              website.Name,
			Short_description: website.Short_description,
			Description:       website.Description,
			Creator_id:        website.Creator_id,
			Banned:            website.Banned,
			Updated_at:        website.Updated_at,
			Created_at:        website.Created_at,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *WebSiteHandler) ListWebSites(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	websites, err := h.webSiteService.ListWebSites(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respData := make([]WebSiteData, 0, len(websites))
	for _, website := range websites {
		respData = append(respData, WebSiteData{
			Id:                website.Id,
			Type:              website.Type,
			Image:             website.Image,
			Name:              website.Name,
			Short_description: website.Short_description,
			Description:       website.Description,
			Creator_id:        website.Creator_id,
			Banned:            website.Banned,
			Updated_at:        website.Updated_at,
			Created_at:        website.Created_at,
		})
	}

	resp := WebSitesResponse{
		Success: true,
		Message: "success",
		Data:    respData,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *WebSiteHandler) ReplaceRoutes(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	siteID := r.PathValue("siteID")
	if strings.TrimSpace(siteID) == "" {
		http.Error(w, "site_id is required", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	var req ReplaceRoutesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	inputs := make([]service.RouteInput, 0, len(req.Routes))
	for _, route := range req.Routes {
		inputs = append(inputs, service.RouteInput{
			Path:         route.Path,
			Title:        route.Title,
			RequiresAuth: route.RequiresAuth,
		})
	}

	routes, err := h.webSiteService.ReplaceRoutes(userID, siteID, inputs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respRoutes := make([]RouteResponse, 0, len(routes))
	for _, route := range routes {
		respRoutes = append(respRoutes, RouteResponse{
			ID:           route.Id,
			Path:         route.Path,
			Title:        route.Title,
			RequiresAuth: route.RequiresAuth,
			Position:     route.Position,
		})
	}

	resp := RoutesResponse{
		Success: true,
		Message: "routes saved",
		Data:    respRoutes,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *WebSiteHandler) ListRoutes(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	siteID := r.PathValue("siteID")
	if strings.TrimSpace(siteID) == "" {
		http.Error(w, "site_id is required", http.StatusBadRequest)
		return
	}

	routes, err := h.webSiteService.ListRoutes(userID, siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respRoutes := make([]RouteResponse, 0, len(routes))
	for _, route := range routes {
		respRoutes = append(respRoutes, RouteResponse{
			ID:           route.Id,
			Path:         route.Path,
			Title:        route.Title,
			RequiresAuth: route.RequiresAuth,
			Position:     route.Position,
		})
	}

	resp := RoutesResponse{
		Success: true,
		Message: "success",
		Data:    respRoutes,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *WebSiteHandler) CreateVersion(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	siteID := r.PathValue("siteID")
	if strings.TrimSpace(siteID) == "" {
		http.Error(w, "site_id is required", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	var req CreateVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	version, scanReport, err := h.webSiteService.CreateVersion(userID, siteID, req.SourceType, req.Source)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := VersionResponse{
		Success: true,
		Message: scanReport.Summary,
		Data: asVersionData(
			version.Id,
			version.WebsiteId,
			version.Version,
			version.SourceType,
			version.Source,
			version.ScanStatus,
			version.ScanScore,
			version.ScanFindings,
			version.Published,
			version.PublishedAt,
		),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *WebSiteHandler) PublishVersion(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	siteID := r.PathValue("siteID")
	if strings.TrimSpace(siteID) == "" {
		http.Error(w, "site_id is required", http.StatusBadRequest)
		return
	}

	versionParam := r.PathValue("version")
	versionNumber, err := strconv.Atoi(versionParam)
	if err != nil || versionNumber <= 0 {
		http.Error(w, "invalid version", http.StatusBadRequest)
		return
	}

	version, err := h.webSiteService.PublishVersion(userID, siteID, versionNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := VersionResponse{
		Success: true,
		Message: "version published",
		Data: asVersionData(
			version.Id,
			version.WebsiteId,
			version.Version,
			version.SourceType,
			version.Source,
			version.ScanStatus,
			version.ScanScore,
			version.ScanFindings,
			version.Published,
			version.PublishedAt,
		),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *WebSiteHandler) GetScanReport(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	siteID := r.PathValue("siteID")
	if strings.TrimSpace(siteID) == "" {
		http.Error(w, "site_id is required", http.StatusBadRequest)
		return
	}

	versionParam := r.PathValue("version")
	versionNumber, err := strconv.Atoi(versionParam)
	if err != nil || versionNumber <= 0 {
		http.Error(w, "invalid version", http.StatusBadRequest)
		return
	}

	version, err := h.webSiteService.GetScanReport(userID, siteID, versionNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := ScanReportResponse{
		Success: true,
		Message: "success",
	}
	resp.Data.Version = version.Version
	resp.Data.ScanStatus = version.ScanStatus
	resp.Data.ScanScore = version.ScanScore
	resp.Data.ScanFindings = version.ScanFindings
	resp.Data.SourceType = version.SourceType

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *WebSiteHandler) ListVersions(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	siteID := r.PathValue("siteID")
	if strings.TrimSpace(siteID) == "" {
		http.Error(w, "site_id is required", http.StatusBadRequest)
		return
	}

	versions, err := h.webSiteService.ListVersions(userID, siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respData := make([]VersionData, 0, len(versions))
	for _, version := range versions {
		respData = append(respData, asVersionData(
			version.Id,
			version.WebsiteId,
			version.Version,
			version.SourceType,
			version.Source,
			version.ScanStatus,
			version.ScanScore,
			version.ScanFindings,
			version.Published,
			version.PublishedAt,
		))
	}

	resp := VersionsResponse{
		Success: true,
		Message: "success",
		Data:    respData,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// PublicStoreInfo is now handled by StoreSocialHandler.GetStoreFullInfo.
// This stub remains for legacy redirect compatibility.
func (h *WebSiteHandler) PublicStoreInfo(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (h *WebSiteHandler) PublicRender(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	if siteID == "" {
		http.NotFound(w, r)
		return
	}

	path := r.PathValue("path")
	if strings.TrimSpace(path) == "" {
		path = "/"
	}

	html, err := h.webSiteService.GetPublicCompiledPage(siteID, path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

func (h *WebSiteHandler) DeleteWebSite(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	siteID := strings.TrimSpace(r.PathValue("siteID"))
	if siteID == "" {
		http.Error(w, "site_id is required", http.StatusBadRequest)
		return
	}

	if err := h.webSiteService.DeleteWebSite(userID, siteID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WebSiteHandler) ListSiteAPIs(w http.ResponseWriter, r *http.Request) {
	resp := SiteAPIsResponse{
		Success: true,
		Message: "success",
		Data: []SiteAPIItem{
			{
				ID:          "store_products",
				Method:      "GET",
				Path:        "/api/store/products",
				Label:       "Listar produtos",
				Description: "Lista produtos públicos da loja.",
			},
			{
				ID:          "store_login",
				Method:      "POST",
				Path:        "/api/store/login",
				Label:       "Login de cliente",
				Description: "Autenticação de clientes do site.",
			},
			{
				ID:          "store_shipping_quote",
				Method:      "POST",
				Path:        "/api/store/shipping/quote",
				Label:       "Calcular frete",
				Description: "Retorna o valor estimado de frete.",
			},
			{
				ID:          "software_download_token",
				Method:      "POST",
				Path:        "/api/software/download-token",
				Label:       "Token de download",
				Description: "Gera token seguro para download de software.",
			},
			{
				ID:          "course_modules",
				Method:      "GET",
				Path:        "/api/course/modules",
				Label:       "Módulos de curso",
				Description: "Lista módulos e aulas publicadas.",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
