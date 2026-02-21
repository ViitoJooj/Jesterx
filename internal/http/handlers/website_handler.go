package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/security"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

type WebSiteRequest struct {
	Type              string `json:"type"`
	Image             []byte `json:"image,omitempty"`
	Name              string `json:"name"`
	Short_description string `json:"short_description,omitempty"`
	Description       string `json:"description,omitempty"`
}

type WebSiteData struct {
	Id                string    `json:"id"`
	Type              string    `json:"type"`
	Image             []byte    `json:"Image"`
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

type WebSiteHandler struct {
	webSiteService *service.WebSiteService
}

func NewWebSiteHandler(webSiteService *service.WebSiteService) *WebSiteHandler {
	return &WebSiteHandler{
		webSiteService: webSiteService,
	}
}

func (h *WebSiteHandler) CreateWebSite(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req WebSiteRequest

	accessCookie, err := r.Cookie("access_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	claims, err := security.ParseAccessToken(accessCookie.Value)
	if err != nil {
		http.Error(w, "Token not valid.", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	website, err := h.webSiteService.CreateWebSite(req.Type, req.Image, req.Name, req.Short_description, req.Description, claims.Sub)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := WebSiteResponse{
		Success: true,
		Message: fmt.Sprintf("%s created.", req.Type),
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
