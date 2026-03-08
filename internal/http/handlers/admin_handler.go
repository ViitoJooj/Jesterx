package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

func (h *AdminHandler) checkAdmin(userID string) bool {
	var role string
	err := h.db.QueryRowContext(context.Background(),
		`SELECT role FROM users WHERE id = $1`, userID).Scan(&role)
	return err == nil && role == "admin"
}

type AdminStats struct {
	TotalUsers    int     `json:"total_users"`
	TotalSites    int     `json:"total_sites"`
	TotalProducts int     `json:"total_products"`
	TotalOrders   int     `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
	PlatformFees  float64 `json:"platform_fees"`
	ActivePlans   int     `json:"active_plans"`
	NewUsersToday int     `json:"new_users_today"`
	OrdersToday   int     `json:"orders_today"`
}

func (h *AdminHandler) Stats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok || !h.checkAdmin(userID) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"success":false,"message":"acesso negado"}`, http.StatusForbidden)
		return
	}

	ctx := context.Background()
	var stats AdminStats
	today := time.Now().Format("2006-01-02")

	h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers)
	h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM websites`).Scan(&stats.TotalSites)
	h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM products`).Scan(&stats.TotalProducts)
	h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders`).Scan(&stats.TotalOrders)
	h.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(total),0) FROM orders WHERE status != 'cancelled'`).Scan(&stats.TotalRevenue)
	h.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(platform_fee),0) FROM orders WHERE status != 'cancelled'`).Scan(&stats.PlatformFees)
	h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM payments WHERE status = 'completed'`).Scan(&stats.ActivePlans)
	h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE DATE(created_at) = $1`, today).Scan(&stats.NewUsersToday)
	h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders WHERE DATE(created_at) = $1`, today).Scan(&stats.OrdersToday)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": stats})
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok || !h.checkAdmin(userID) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"success":false,"message":"acesso negado"}`, http.StatusForbidden)
		return
	}

	ctx := context.Background()
	rows, err := h.db.QueryContext(ctx, `
		SELECT u.id, u.first_name, u.last_name, u.email, u.role, u.verified_email,
		       u.created_at,
		       COUNT(DISTINCT w.id) as site_count
		FROM users u
		LEFT JOIN websites w ON w.creator_id = u.id
		GROUP BY u.id
		ORDER BY u.created_at DESC
		LIMIT 100
	`)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type UserRow struct {
		ID            string    `json:"id"`
		FirstName     string    `json:"first_name"`
		LastName      string    `json:"last_name"`
		Email         string    `json:"email"`
		Role          string    `json:"role"`
		VerifiedEmail bool      `json:"verified_email"`
		CreatedAt     time.Time `json:"created_at"`
		SiteCount     int       `json:"site_count"`
	}

	users := make([]UserRow, 0)
	for rows.Next() {
		var u UserRow
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role,
			&u.VerifiedEmail, &u.CreatedAt, &u.SiteCount); err != nil {
			continue
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": users})
}

func (h *AdminHandler) ListSites(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok || !h.checkAdmin(userID) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"success":false,"message":"acesso negado"}`, http.StatusForbidden)
		return
	}

	ctx := context.Background()
	rows, err := h.db.QueryContext(ctx, `
		SELECT w.id, w.name, w.website_type, w.banned, w.created_at,
		       u.first_name || ' ' || u.last_name as owner_name, u.email as owner_email,
		       COUNT(DISTINCT wv.id) as version_count
		FROM websites w
		LEFT JOIN users u ON u.id = w.creator_id
		LEFT JOIN website_versions wv ON wv.website_id = w.id
		GROUP BY w.id, u.first_name, u.last_name, u.email
		ORDER BY w.created_at DESC
		LIMIT 200
	`)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type SiteRow struct {
		ID           string    `json:"id"`
		Name         string    `json:"name"`
		Type         string    `json:"type"`
		Banned       bool      `json:"banned"`
		CreatedAt    time.Time `json:"created_at"`
		OwnerName    string    `json:"owner_name"`
		OwnerEmail   string    `json:"owner_email"`
		VersionCount int       `json:"version_count"`
	}

	sites := make([]SiteRow, 0)
	for rows.Next() {
		var s SiteRow
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.Banned, &s.CreatedAt,
			&s.OwnerName, &s.OwnerEmail, &s.VersionCount); err != nil {
			continue
		}
		sites = append(sites, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": sites})
}

func (h *AdminHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok || !h.checkAdmin(userID) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"success":false,"message":"acesso negado"}`, http.StatusForbidden)
		return
	}

	ctx := context.Background()
	rows, err := h.db.QueryContext(ctx, `
		SELECT o.id, o.website_id, COALESCE(w.name, '') as site_name,
		       o.buyer_name, o.buyer_email, o.status,
		       o.subtotal, o.platform_fee, o.total, o.created_at
		FROM orders o
		LEFT JOIN websites w ON w.id = o.website_id
		ORDER BY o.created_at DESC
		LIMIT 500
	`)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type OrderRow struct {
		ID          string    `json:"id"`
		WebsiteID   string    `json:"website_id"`
		SiteName    string    `json:"site_name"`
		BuyerName   string    `json:"buyer_name"`
		BuyerEmail  string    `json:"buyer_email"`
		Status      string    `json:"status"`
		Subtotal    float64   `json:"subtotal"`
		PlatformFee float64   `json:"platform_fee"`
		Total       float64   `json:"total"`
		CreatedAt   time.Time `json:"created_at"`
	}

	orders := make([]OrderRow, 0)
	for rows.Next() {
		var o OrderRow
		if err := rows.Scan(&o.ID, &o.WebsiteID, &o.SiteName,
			&o.BuyerName, &o.BuyerEmail, &o.Status,
			&o.Subtotal, &o.PlatformFee, &o.Total, &o.CreatedAt); err != nil {
			continue
		}
		orders = append(orders, o)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": orders})
}

func (h *AdminHandler) Revenue(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok || !h.checkAdmin(userID) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"success":false,"message":"acesso negado"}`, http.StatusForbidden)
		return
	}

	ctx := context.Background()
	rows, err := h.db.QueryContext(ctx, `
		SELECT DATE(created_at) as day,
		       COUNT(*) as order_count,
		       COALESCE(SUM(total), 0) as gmv,
		       COALESCE(SUM(platform_fee), 0) as fee
		FROM orders
		WHERE created_at >= NOW() - INTERVAL '30 days'
		  AND status != 'cancelled'
		GROUP BY DATE(created_at)
		ORDER BY day ASC
	`)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type DayRevenue struct {
		Day        string  `json:"day"`
		OrderCount int     `json:"order_count"`
		GMV        float64 `json:"gmv"`
		Fee        float64 `json:"fee"`
	}

	days := make([]DayRevenue, 0)
	for rows.Next() {
		var d DayRevenue
		var day time.Time
		if err := rows.Scan(&day, &d.OrderCount, &d.GMV, &d.Fee); err != nil {
			continue
		}
		d.Day = day.Format("2006-01-02")
		days = append(days, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": days})
}
