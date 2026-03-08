package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
)

type ThemeHandler struct {
	db *sql.DB
}

func NewThemeHandler(db *sql.DB) *ThemeHandler {
	return &ThemeHandler{db: db}
}

type ThemeData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	PreviewURL  string `json:"preview_url"`
	SourceType  string `json:"source_type"`
	Source      string `json:"source"`
}

func (h *ThemeHandler) ListThemes(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(context.Background(),
		`SELECT id, name, COALESCE(description,''), category, COALESCE(preview_url,''), source_type, source
		 FROM themes WHERE active = true ORDER BY name`)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	themes := make([]ThemeData, 0)
	for rows.Next() {
		var t ThemeData
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Category, &t.PreviewURL, &t.SourceType, &t.Source); err != nil {
			continue
		}
		themes = append(themes, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "success", "data": themes})
}
