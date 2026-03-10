package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

type StoreSocialHandler struct {
	svc *service.StoreSocialService
	db  *sql.DB // used only for admin role check
}

func NewStoreSocialHandler(svc *service.StoreSocialService, db *sql.DB) *StoreSocialHandler {
	return &StoreSocialHandler{svc: svc, db: db}
}

func (h *StoreSocialHandler) checkAdmin(userID string) bool {
	var role string
	err := h.db.QueryRowContext(context.Background(),
		`SELECT role FROM users WHERE id = $1`, userID).Scan(&role)
	return err == nil && role == "admin"
}

// ─── Store Full Info ──────────────────────────────────────────────────────────

// GET /api/store/{siteID}/info
func (h *StoreSocialHandler) GetStoreFullInfo(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	if siteID == "" {
		http.NotFound(w, r)
		return
	}

	info, err := h.svc.GetStoreFullInfo(siteID)
	if err != nil {
		http.Error(w, "loja não encontrada", http.StatusNotFound)
		return
	}

	// Record visit asynchronously
	go func() { _ = h.svc.RecordVisit(siteID) }()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "success", "data": info})
}

// ─── Visit Stats ─────────────────────────────────────────────────────────────

// GET /api/store/{siteID}/visits?days=30
func (h *StoreSocialHandler) GetVisitStats(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days <= 0 {
		days = 30
	}

	stats, err := h.svc.GetVisitStats(siteID, days)
	if err != nil {
		http.Error(w, "erro ao buscar estatísticas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": stats})
}

// ─── Comments ────────────────────────────────────────────────────────────────

type CommentData struct {
	ID              string        `json:"id"`
	WebsiteID       string        `json:"website_id"`
	UserID          string        `json:"user_id"`
	UserName        string        `json:"user_name"`
	AvatarURL       *string       `json:"avatar_url,omitempty"`
	Content         string        `json:"content"`
	ParentCommentID *string       `json:"parent_comment_id,omitempty"`
	Replies         []CommentData `json:"replies,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
}

func commentToData(c *domain.StoreComment) CommentData {
	d := CommentData{
		ID: c.ID, WebsiteID: c.WebsiteID, UserID: c.UserID,
		UserName: c.UserName, AvatarURL: c.AvatarURL,
		Content: c.Content, ParentCommentID: c.ParentCommentID,
		CreatedAt: c.CreatedAt,
	}
	for i := range c.Replies {
		d.Replies = append(d.Replies, commentToData(&c.Replies[i]))
	}
	return d
}

// GET /api/store/{siteID}/comments
func (h *StoreSocialHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	comments, err := h.svc.ListComments(siteID)
	if err != nil {
		http.Error(w, "erro ao buscar comentários", http.StatusInternalServerError)
		return
	}
	data := make([]CommentData, 0, len(comments))
	for i := range comments {
		data = append(data, commentToData(&comments[i]))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": data})
}

// POST /api/store/{siteID}/comments  (auth required)
func (h *StoreSocialHandler) PostComment(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Content string `json:"content"`
		Stars   int    `json:"stars"`
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	comment, err := h.svc.PostComment(siteID, userID, strings.TrimSpace(req.Content), req.Stars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": commentToData(comment)})
}

// DELETE /api/store/{siteID}/comments/{commentID}  (auth required)
func (h *StoreSocialHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	commentID := strings.TrimSpace(r.PathValue("commentID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.svc.DeleteComment(commentID, userID, siteID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "comentário removido"})
}

// POST /api/store/{siteID}/comments/{commentID}/replies  (auth + team role required)
func (h *StoreSocialHandler) ReplyComment(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	commentID := strings.TrimSpace(r.PathValue("commentID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	reply, err := h.svc.ReplyComment(siteID, commentID, userID, strings.TrimSpace(req.Content))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "sem permissão para responder comentários" {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": commentToData(reply)})
}

// ─── Ratings ─────────────────────────────────────────────────────────────────

// POST /api/store/{siteID}/ratings  (auth required)
func (h *StoreSocialHandler) RateStore(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Stars int `json:"stars"`
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	rating, err := h.svc.RateStore(siteID, userID, req.Stars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": rating})
}

// GET /api/store/{siteID}/my-rating  (auth required)
func (h *StoreSocialHandler) GetMyRating(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rating, err := h.svc.GetMyRating(siteID, userID)
	if err != nil {
		http.Error(w, "erro ao buscar avaliação", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": rating})
}

// ─── Owner: Update Store Profile ─────────────────────────────────────────────

// PATCH /api/v1/sites/{siteID}/profile  (auth + owner required)
func (h *StoreSocialHandler) UpdateStoreProfile(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name             string `json:"name"`
		ShortDescription string `json:"short_description"`
		Description      string `json:"description"`
		Image            []byte `json:"image,omitempty"`
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateStoreProfile(siteID, userID,
		strings.TrimSpace(req.Name),
		strings.TrimSpace(req.ShortDescription),
		strings.TrimSpace(req.Description),
		req.Image,
	); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "perfil atualizado"})
}

// ─── Admin: Set Mature Content ────────────────────────────────────────────────

// PATCH /api/v1/admin/sites/{siteID}/mature  (admin only)
func (h *StoreSocialHandler) AdminSetMature(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok || !h.checkAdmin(userID) {
		http.Error(w, "acesso negado", http.StatusForbidden)
		return
	}

	var req struct {
		MatureContent bool `json:"mature_content"`
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.svc.SetMatureContent(siteID, req.MatureContent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": map[bool]string{true: "loja marcada como +18", false: "tag +18 removida"}[req.MatureContent],
	})
}

// ─── Team Members ─────────────────────────────────────────────────────────────

// GET /api/v1/sites/{siteID}/members  (auth + any role)
func (h *StoreSocialHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	members, err := h.svc.ListMembers(siteID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": members})
}

// POST /api/v1/sites/{siteID}/members  (auth + owner/manager/admin)
func (h *StoreSocialHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.Role == "" {
		http.Error(w, "user_id e role são obrigatórios", http.StatusBadRequest)
		return
	}

	member, err := h.svc.AddMember(siteID, userID, strings.TrimSpace(req.UserID), strings.TrimSpace(req.Role))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": member})
}

// DELETE /api/v1/sites/{siteID}/members/{memberUserID}  (auth + owner/manager/admin)
func (h *StoreSocialHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	targetUserID := strings.TrimSpace(r.PathValue("memberUserID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.svc.RemoveMember(siteID, userID, targetUserID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "membro removido"})
}

// PATCH /api/v1/sites/{siteID}/members/{memberUserID}  (auth + owner/manager/admin)
func (h *StoreSocialHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	targetUserID := strings.TrimSpace(r.PathValue("memberUserID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Role string `json:"role"`
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Role) == "" {
		http.Error(w, "role é obrigatória", http.StatusBadRequest)
		return
	}

	member, err := h.svc.UpdateMemberRole(siteID, userID, targetUserID, strings.TrimSpace(req.Role))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": member})
}

// ─── My Role ──────────────────────────────────────────────────────────────────

// GET /api/store/{siteID}/my-role  (auth required)
func (h *StoreSocialHandler) GetMyRole(w http.ResponseWriter, r *http.Request) {
	siteID := strings.TrimSpace(r.PathValue("siteID"))
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	role, err := h.svc.GetUserRoleInStore(userID, siteID)
	if err != nil {
		// Store not found — user has no role
		role = ""
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": map[string]string{"role": role}})
}
